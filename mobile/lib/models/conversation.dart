import 'user.dart';

class QubeConversation {
  final String id;
  final bool isGroup;
  final String name;
  final List<QubeUser> participants;
  final QubeMessage? lastMessage;
  final int unreadCount;
  final DateTime updatedAt;

  QubeConversation({
    required this.id,
    this.isGroup = false,
    this.name = '',
    this.participants = const [],
    this.lastMessage,
    this.unreadCount = 0,
    required this.updatedAt,
  });

  factory QubeConversation.fromJson(Map<String, dynamic> json) {
    return QubeConversation(
      id: json['id'] as String,
      isGroup: json['isGroup'] as bool? ?? false,
      name: json['name'] as String? ?? '',
      participants: (json['participants'] as List?)?.map((p) => QubeUser.fromJson(p)).toList() ?? [],
      lastMessage: json['lastMessage'] != null ? QubeMessage.fromJson(json['lastMessage']) : null,
      unreadCount: json['unreadCount'] as int? ?? 0,
      updatedAt: DateTime.parse(json['updatedAt'] as String),
    );
  }
}

class QubeMessage {
  final String id;
  final String conversationId;
  final QubeUser sender;
  final String content;
  final DateTime createdAt;

  QubeMessage({
    required this.id,
    required this.conversationId,
    required this.sender,
    required this.content,
    required this.createdAt,
  });

  factory QubeMessage.fromJson(Map<String, dynamic> json) {
    return QubeMessage(
      id: json['id'] as String,
      conversationId: json['conversationId'] as String,
      sender: QubeUser.fromJson(json['sender'] as Map<String, dynamic>),
      content: json['content'] as String,
      createdAt: DateTime.parse(json['createdAt'] as String),
    );
  }
}
