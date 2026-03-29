import 'package:flutter/material.dart';
import 'package:timeago/timeago.dart' as timeago;
import '../../config/theme.dart';
import '../../models/conversation.dart';
import '../../services/api_service.dart';
import 'chat_screen.dart';

class ConversationsScreen extends StatefulWidget {
  const ConversationsScreen({super.key});

  @override
  State<ConversationsScreen> createState() => _ConversationsScreenState();
}

class _ConversationsScreenState extends State<ConversationsScreen> {
  final _api = ApiService();
  List<QubeConversation> _conversations = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _loadConversations();
  }

  Future<void> _loadConversations() async {
    try {
      final data = await _api.query('Conversations', variables: {'limit': 20});
      setState(() {
        _conversations = (data['conversations']['conversations'] as List)
            .map((c) => QubeConversation.fromJson(c))
            .toList();
        _isLoading = false;
      });
    } catch (e) {
      setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Messages'),
        actions: [
          IconButton(
            icon: const Icon(Icons.edit_outlined),
            onPressed: () {
              // TODO: New conversation
            },
          ),
        ],
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : _conversations.isEmpty
              ? const Center(
                  child: Text('No messages yet', style: TextStyle(color: QubeTheme.textSecondary)),
                )
              : RefreshIndicator(
                  onRefresh: _loadConversations,
                  child: ListView.builder(
                    itemCount: _conversations.length,
                    itemBuilder: (context, index) {
                      final conv = _conversations[index];
                      final otherUser = conv.participants.isNotEmpty ? conv.participants.first : null;

                      return ListTile(
                        leading: CircleAvatar(
                          backgroundColor: QubeTheme.surface,
                          backgroundImage: otherUser?.avatarUrl.isNotEmpty == true
                              ? NetworkImage(otherUser!.avatarUrl)
                              : null,
                          child: otherUser?.avatarUrl.isEmpty != false
                              ? Text(otherUser?.displayName[0].toUpperCase() ?? '?')
                              : null,
                        ),
                        title: Text(
                          conv.isGroup ? conv.name : (otherUser?.displayName ?? 'Unknown'),
                          style: const TextStyle(fontWeight: FontWeight.w600),
                        ),
                        subtitle: conv.lastMessage != null
                            ? Text(
                                conv.lastMessage!.content,
                                maxLines: 1,
                                overflow: TextOverflow.ellipsis,
                                style: TextStyle(
                                  color: conv.unreadCount > 0 ? QubeTheme.textPrimary : QubeTheme.textSecondary,
                                ),
                              )
                            : null,
                        trailing: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Text(
                              timeago.format(conv.updatedAt, locale: 'ja'),
                              style: const TextStyle(color: QubeTheme.textSecondary, fontSize: 12),
                            ),
                            if (conv.unreadCount > 0) ...[
                              const SizedBox(height: 4),
                              Container(
                                padding: const EdgeInsets.all(6),
                                decoration: const BoxDecoration(
                                  color: QubeTheme.primary,
                                  shape: BoxShape.circle,
                                ),
                                child: Text('${conv.unreadCount}', style: const TextStyle(fontSize: 11)),
                              ),
                            ],
                          ],
                        ),
                        onTap: () {
                          Navigator.push(
                            context,
                            MaterialPageRoute(
                              builder: (_) => ChatScreen(conversationId: conv.id, title: otherUser?.displayName ?? conv.name),
                            ),
                          );
                        },
                      );
                    },
                  ),
                ),
    );
  }
}
