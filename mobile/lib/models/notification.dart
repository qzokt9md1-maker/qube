import 'user.dart';

class QubeNotification {
  final String id;
  final QubeUser actor;
  final String type;
  final String? postId;
  final bool isRead;
  final DateTime createdAt;

  QubeNotification({
    required this.id,
    required this.actor,
    required this.type,
    this.postId,
    this.isRead = false,
    required this.createdAt,
  });

  factory QubeNotification.fromJson(Map<String, dynamic> json) {
    return QubeNotification(
      id: json['id'] as String,
      actor: QubeUser.fromJson(json['actor'] as Map<String, dynamic>),
      type: json['type'] as String,
      postId: json['postId'] as String?,
      isRead: json['isRead'] as bool? ?? false,
      createdAt: DateTime.parse(json['createdAt'] as String),
    );
  }

  String get displayText {
    switch (type) {
      case 'like':
        return '${actor.displayName}があなたの投稿にいいねしました';
      case 'repost':
        return '${actor.displayName}があなたの投稿をリポストしました';
      case 'follow':
        return '${actor.displayName}があなたをフォローしました';
      case 'reply':
        return '${actor.displayName}があなたの投稿に返信しました';
      case 'quote':
        return '${actor.displayName}があなたの投稿を引用しました';
      case 'mention':
        return '${actor.displayName}があなたをメンションしました';
      case 'dm':
        return '${actor.displayName}からメッセージが届きました';
      default:
        return '${actor.displayName}からの通知';
    }
  }
}
