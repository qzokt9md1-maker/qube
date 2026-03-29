import 'dart:async';
import 'package:flutter/material.dart';
import '../../config/theme.dart';
import '../../models/conversation.dart';
import '../../services/api_service.dart';
import '../../services/ws_service.dart';

class ChatScreen extends StatefulWidget {
  final String conversationId;
  final String title;

  const ChatScreen({super.key, required this.conversationId, required this.title});

  @override
  State<ChatScreen> createState() => _ChatScreenState();
}

class _ChatScreenState extends State<ChatScreen> {
  final _api = ApiService();
  final _ws = WsService();
  final _controller = TextEditingController();
  final _scrollController = ScrollController();
  List<QubeMessage> _messages = [];
  bool _isLoading = true;
  StreamSubscription? _wsSub;

  @override
  void initState() {
    super.initState();
    _loadMessages();
    _markRead();
    _wsSub = _ws.stream.listen((event) {
      if (event['type'] == 'new_message') {
        final payload = event['payload'] as Map<String, dynamic>;
        if (payload['conversation_id'] == widget.conversationId) {
          _loadMessages();
        }
      }
    });
  }

  @override
  void dispose() {
    _controller.dispose();
    _scrollController.dispose();
    _wsSub?.cancel();
    super.dispose();
  }

  Future<void> _loadMessages() async {
    try {
      final data = await _api.query('Messages', variables: {
        'conversationId': widget.conversationId,
        'limit': 30,
      });
      setState(() {
        _messages = (data['messages']['messages'] as List)
            .map((m) => QubeMessage.fromJson(m))
            .toList()
            .reversed
            .toList();
        _isLoading = false;
      });
      _scrollToBottom();
    } catch (e) {
      setState(() => _isLoading = false);
    }
  }

  Future<void> _markRead() async {
    await _api.query('MarkConversationRead', variables: {'conversationId': widget.conversationId});
  }

  Future<void> _sendMessage() async {
    final text = _controller.text.trim();
    if (text.isEmpty) return;
    _controller.clear();
    try {
      await _api.query('SendMessage', variables: {
        'input': {
          'conversationId': widget.conversationId,
          'content': text,
        },
      });
      _loadMessages();
    } catch (e) {
      _controller.text = text;
    }
  }

  void _scrollToBottom() {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (_scrollController.hasClients) {
        _scrollController.animateTo(
          _scrollController.position.maxScrollExtent,
          duration: const Duration(milliseconds: 200),
          curve: Curves.easeOut,
        );
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text(widget.title)),
      body: Column(
        children: [
          Expanded(
            child: _isLoading
                ? const Center(child: CircularProgressIndicator())
                : ListView.builder(
                    controller: _scrollController,
                    padding: const EdgeInsets.all(16),
                    itemCount: _messages.length,
                    itemBuilder: (context, index) {
                      final msg = _messages[index];
                      return _MessageBubble(message: msg);
                    },
                  ),
          ),
          Container(
            padding: const EdgeInsets.all(8),
            decoration: const BoxDecoration(
              color: QubeTheme.surface,
              border: Border(top: BorderSide(color: QubeTheme.border, width: 0.5)),
            ),
            child: SafeArea(
              child: Row(
                children: [
                  Expanded(
                    child: TextField(
                      controller: _controller,
                      decoration: InputDecoration(
                        hintText: 'Message...',
                        filled: true,
                        fillColor: QubeTheme.background,
                        border: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(24),
                          borderSide: BorderSide.none,
                        ),
                        contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
                      ),
                      onChanged: (_) {
                        _ws.sendTyping(widget.conversationId);
                      },
                      onSubmitted: (_) => _sendMessage(),
                    ),
                  ),
                  const SizedBox(width: 8),
                  IconButton(
                    icon: const Icon(Icons.send, color: QubeTheme.primary),
                    onPressed: _sendMessage,
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _MessageBubble extends StatelessWidget {
  final QubeMessage message;
  const _MessageBubble({required this.message});

  @override
  Widget build(BuildContext context) {
    // TODO: Determine if sent by current user for alignment
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          CircleAvatar(
            radius: 16,
            backgroundColor: QubeTheme.surface,
            child: Text(message.sender.displayName[0].toUpperCase(), style: const TextStyle(fontSize: 12)),
          ),
          const SizedBox(width: 8),
          Flexible(
            child: Container(
              padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
              decoration: BoxDecoration(
                color: QubeTheme.surface,
                borderRadius: BorderRadius.circular(16),
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    message.sender.displayName,
                    style: const TextStyle(fontWeight: FontWeight.w600, fontSize: 13),
                  ),
                  const SizedBox(height: 2),
                  Text(message.content, style: const TextStyle(fontSize: 15)),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}
