import 'package:flutter/material.dart';
import 'package:timeago/timeago.dart' as timeago;
import '../../config/theme.dart';
import '../../models/notification.dart';
import '../../services/api_service.dart';

class NotificationsScreen extends StatefulWidget {
  const NotificationsScreen({super.key});

  @override
  State<NotificationsScreen> createState() => _NotificationsScreenState();
}

class _NotificationsScreenState extends State<NotificationsScreen> {
  final _api = ApiService();
  List<QubeNotification> _notifications = [];
  bool _isLoading = true;
  int _unreadCount = 0;

  @override
  void initState() {
    super.initState();
    _loadNotifications();
  }

  Future<void> _loadNotifications() async {
    try {
      final data = await _api.query('Notifications', variables: {'limit': 30});
      final result = data['notifications'];
      setState(() {
        _notifications = (result['notifications'] as List)
            .map((n) => QubeNotification.fromJson(n))
            .toList();
        _unreadCount = result['unreadCount'] as int? ?? 0;
        _isLoading = false;
      });
    } catch (e) {
      setState(() => _isLoading = false);
    }
  }

  Future<void> _markAllRead() async {
    await _api.query('MarkAllNotificationsRead');
    setState(() {
      _unreadCount = 0;
    });
  }

  IconData _iconForType(String type) {
    switch (type) {
      case 'like':
        return Icons.favorite;
      case 'repost':
        return Icons.repeat;
      case 'follow':
        return Icons.person_add;
      case 'reply':
        return Icons.reply;
      case 'quote':
        return Icons.format_quote;
      case 'mention':
        return Icons.alternate_email;
      case 'dm':
        return Icons.mail;
      default:
        return Icons.notifications;
    }
  }

  Color _colorForType(String type) {
    switch (type) {
      case 'like':
        return QubeTheme.danger;
      case 'repost':
        return QubeTheme.success;
      case 'follow':
        return QubeTheme.primary;
      default:
        return QubeTheme.textSecondary;
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Notifications'),
        actions: [
          if (_unreadCount > 0)
            TextButton(
              onPressed: _markAllRead,
              child: const Text('Mark all read'),
            ),
        ],
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : _notifications.isEmpty
              ? const Center(
                  child: Text('No notifications', style: TextStyle(color: QubeTheme.textSecondary)),
                )
              : RefreshIndicator(
                  onRefresh: _loadNotifications,
                  child: ListView.builder(
                    itemCount: _notifications.length,
                    itemBuilder: (context, index) {
                      final notif = _notifications[index];
                      return Container(
                        color: notif.isRead ? null : QubeTheme.primary.withValues(alpha: 0.05),
                        child: ListTile(
                          leading: CircleAvatar(
                            backgroundColor: _colorForType(notif.type).withValues(alpha: 0.15),
                            child: Icon(_iconForType(notif.type), color: _colorForType(notif.type), size: 20),
                          ),
                          title: Text(notif.displayText, style: const TextStyle(fontSize: 14)),
                          subtitle: Text(
                            timeago.format(notif.createdAt, locale: 'ja'),
                            style: const TextStyle(color: QubeTheme.textSecondary, fontSize: 12),
                          ),
                          onTap: () {
                            // TODO: Navigate to relevant content
                          },
                        ),
                      );
                    },
                  ),
                ),
    );
  }
}
