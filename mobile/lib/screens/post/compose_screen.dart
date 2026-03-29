import 'package:flutter/material.dart';
import '../../config/theme.dart';
import '../../config/constants.dart';
import '../../services/api_service.dart';

class ComposeScreen extends StatefulWidget {
  final String? replyToId;
  final String? quoteOfId;

  const ComposeScreen({super.key, this.replyToId, this.quoteOfId});

  @override
  State<ComposeScreen> createState() => _ComposeScreenState();
}

class _ComposeScreenState extends State<ComposeScreen> {
  final _controller = TextEditingController();
  final _api = ApiService();
  bool _isPosting = false;

  int get _remaining => QubeConstants.maxPostLength - _controller.text.length;

  Future<void> _post() async {
    if (_controller.text.trim().isEmpty || _isPosting) return;
    setState(() => _isPosting = true);
    try {
      await _api.query('CreatePost', variables: {
        'input': {
          'content': _controller.text.trim(),
          if (widget.replyToId != null) 'replyToId': widget.replyToId,
          if (widget.quoteOfId != null) 'quoteOfId': widget.quoteOfId,
        },
      });
      if (mounted) Navigator.pop(context, true);
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error: $e')),
        );
      }
      setState(() => _isPosting = false);
    }
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        leading: IconButton(
          icon: const Icon(Icons.close),
          onPressed: () => Navigator.pop(context),
        ),
        actions: [
          Padding(
            padding: const EdgeInsets.only(right: 12),
            child: ElevatedButton(
              onPressed: _controller.text.trim().isEmpty || _isPosting ? null : _post,
              style: ElevatedButton.styleFrom(
                padding: const EdgeInsets.symmetric(horizontal: 20),
              ),
              child: _isPosting
                  ? const SizedBox(height: 16, width: 16, child: CircularProgressIndicator(strokeWidth: 2))
                  : const Text('Post'),
            ),
          ),
        ],
      ),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            Expanded(
              child: TextField(
                controller: _controller,
                maxLength: QubeConstants.maxPostLength,
                maxLines: null,
                autofocus: true,
                decoration: InputDecoration(
                  hintText: widget.replyToId != null ? 'Post your reply...' : "What's happening?",
                  border: InputBorder.none,
                  counterText: '',
                ),
                style: const TextStyle(fontSize: 18),
                onChanged: (_) => setState(() {}),
              ),
            ),
            const Divider(color: QubeTheme.border),
            Row(
              children: [
                IconButton(
                  icon: const Icon(Icons.image_outlined, color: QubeTheme.primary),
                  onPressed: () {
                    // TODO: Image picker
                  },
                ),
                IconButton(
                  icon: const Icon(Icons.gif_box_outlined, color: QubeTheme.primary),
                  onPressed: () {},
                ),
                const Spacer(),
                Text(
                  '$_remaining',
                  style: TextStyle(
                    color: _remaining < 0
                        ? QubeTheme.danger
                        : _remaining < 20
                            ? Colors.orange
                            : QubeTheme.textSecondary,
                    fontSize: 14,
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
