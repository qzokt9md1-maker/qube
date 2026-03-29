class QubeUser {
  final String id;
  final String username;
  final String displayName;
  final String bio;
  final String avatarUrl;
  final String headerUrl;
  final String location;
  final String website;
  final bool isVerified;
  final bool isPrivate;
  final int followerCount;
  final int followingCount;
  final int postCount;
  final DateTime createdAt;

  QubeUser({
    required this.id,
    required this.username,
    required this.displayName,
    this.bio = '',
    this.avatarUrl = '',
    this.headerUrl = '',
    this.location = '',
    this.website = '',
    this.isVerified = false,
    this.isPrivate = false,
    this.followerCount = 0,
    this.followingCount = 0,
    this.postCount = 0,
    required this.createdAt,
  });

  factory QubeUser.fromJson(Map<String, dynamic> json) {
    return QubeUser(
      id: json['id'] as String,
      username: json['username'] as String,
      displayName: json['displayName'] as String,
      bio: json['bio'] as String? ?? '',
      avatarUrl: json['avatarUrl'] as String? ?? '',
      headerUrl: json['headerUrl'] as String? ?? '',
      location: json['location'] as String? ?? '',
      website: json['website'] as String? ?? '',
      isVerified: json['isVerified'] as bool? ?? false,
      isPrivate: json['isPrivate'] as bool? ?? false,
      followerCount: json['followerCount'] as int? ?? 0,
      followingCount: json['followingCount'] as int? ?? 0,
      postCount: json['postCount'] as int? ?? 0,
      createdAt: DateTime.parse(json['createdAt'] as String),
    );
  }
}
