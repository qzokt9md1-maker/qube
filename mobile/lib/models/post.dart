import 'user.dart';

class QubePost {
  final String id;
  final QubeUser user;
  final String content;
  final List<QubeMedia> media;
  final String? replyToId;
  final String? repostOfId;
  final String? quoteOfId;
  final int likeCount;
  final int repostCount;
  final int replyCount;
  final int quoteCount;
  final bool isLiked;
  final bool isReposted;
  final bool isBookmarked;
  final DateTime createdAt;

  QubePost({
    required this.id,
    required this.user,
    required this.content,
    this.media = const [],
    this.replyToId,
    this.repostOfId,
    this.quoteOfId,
    this.likeCount = 0,
    this.repostCount = 0,
    this.replyCount = 0,
    this.quoteCount = 0,
    this.isLiked = false,
    this.isReposted = false,
    this.isBookmarked = false,
    required this.createdAt,
  });

  factory QubePost.fromJson(Map<String, dynamic> json) {
    return QubePost(
      id: json['id'] as String,
      user: QubeUser.fromJson(json['user'] as Map<String, dynamic>),
      content: json['content'] as String? ?? '',
      media: (json['media'] as List?)?.map((m) => QubeMedia.fromJson(m)).toList() ?? [],
      replyToId: json['replyToId'] as String?,
      repostOfId: json['repostOfId'] as String?,
      quoteOfId: json['quoteOfId'] as String?,
      likeCount: json['likeCount'] as int? ?? 0,
      repostCount: json['repostCount'] as int? ?? 0,
      replyCount: json['replyCount'] as int? ?? 0,
      quoteCount: json['quoteCount'] as int? ?? 0,
      isLiked: json['isLiked'] as bool? ?? false,
      isReposted: json['isReposted'] as bool? ?? false,
      isBookmarked: json['isBookmarked'] as bool? ?? false,
      createdAt: DateTime.parse(json['createdAt'] as String),
    );
  }
}

class QubeMedia {
  final String id;
  final String mediaType;
  final String url;
  final String thumbnailUrl;
  final int? width;
  final int? height;

  QubeMedia({
    required this.id,
    required this.mediaType,
    required this.url,
    this.thumbnailUrl = '',
    this.width,
    this.height,
  });

  factory QubeMedia.fromJson(Map<String, dynamic> json) {
    return QubeMedia(
      id: json['id'] as String,
      mediaType: json['mediaType'] as String,
      url: json['url'] as String,
      thumbnailUrl: json['thumbnailUrl'] as String? ?? '',
      width: json['width'] as int?,
      height: json['height'] as int?,
    );
  }
}
