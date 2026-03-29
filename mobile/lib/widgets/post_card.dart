import 'package:flutter/material.dart';
import 'package:timeago/timeago.dart' as timeago;
import '../config/theme.dart';
import '../models/post.dart';

class PostCard extends StatelessWidget {
  final QubePost post;
  final VoidCallback? onTap;
  final VoidCallback? onLike;
  final VoidCallback? onRepost;
  final VoidCallback? onReply;
  final VoidCallback? onBookmark;
  final VoidCallback? onUserTap;

  const PostCard({
    super.key,
    required this.post,
    this.onTap,
    this.onLike,
    this.onRepost,
    this.onReply,
    this.onBookmark,
    this.onUserTap,
  });

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        decoration: const BoxDecoration(
          border: Border(bottom: BorderSide(color: QubeTheme.border, width: 0.5)),
        ),
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            GestureDetector(
              onTap: onUserTap,
              child: CircleAvatar(
                radius: 22,
                backgroundColor: QubeTheme.surface,
                backgroundImage: post.user.avatarUrl.isNotEmpty
                    ? NetworkImage(post.user.avatarUrl)
                    : null,
                child: post.user.avatarUrl.isEmpty
                    ? Text(
                        post.user.displayName[0].toUpperCase(),
                        style: const TextStyle(color: QubeTheme.textPrimary, fontSize: 18),
                      )
                    : null,
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Header
                  Row(
                    children: [
                      Flexible(
                        child: Text(
                          post.user.displayName,
                          style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 15),
                          overflow: TextOverflow.ellipsis,
                        ),
                      ),
                      if (post.user.isVerified) ...[
                        const SizedBox(width: 4),
                        const Icon(Icons.verified, size: 16, color: QubeTheme.primary),
                      ],
                      const SizedBox(width: 4),
                      Text(
                        '@${post.user.username}',
                        style: const TextStyle(color: QubeTheme.textSecondary, fontSize: 14),
                      ),
                      const SizedBox(width: 4),
                      Text(
                        '· ${timeago.format(post.createdAt, locale: 'ja')}',
                        style: const TextStyle(color: QubeTheme.textSecondary, fontSize: 14),
                      ),
                    ],
                  ),
                  const SizedBox(height: 4),
                  // Content
                  Text(
                    post.content,
                    style: const TextStyle(fontSize: 15, height: 1.4),
                  ),
                  // Media
                  if (post.media.isNotEmpty) ...[
                    const SizedBox(height: 8),
                    ClipRRect(
                      borderRadius: BorderRadius.circular(12),
                      child: post.media.length == 1
                          ? Image.network(
                              post.media[0].url,
                              fit: BoxFit.cover,
                              width: double.infinity,
                              height: 200,
                            )
                          : SizedBox(
                              height: 200,
                              child: ListView.separated(
                                scrollDirection: Axis.horizontal,
                                itemCount: post.media.length,
                                separatorBuilder: (_, __) => const SizedBox(width: 4),
                                itemBuilder: (context, index) {
                                  return ClipRRect(
                                    borderRadius: BorderRadius.circular(8),
                                    child: Image.network(
                                      post.media[index].url,
                                      fit: BoxFit.cover,
                                      width: 160,
                                    ),
                                  );
                                },
                              ),
                            ),
                    ),
                  ],
                  // Actions
                  const SizedBox(height: 8),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      _ActionButton(
                        icon: Icons.chat_bubble_outline,
                        count: post.replyCount,
                        onTap: onReply,
                      ),
                      _ActionButton(
                        icon: Icons.repeat,
                        count: post.repostCount,
                        isActive: post.isReposted,
                        activeColor: QubeTheme.success,
                        onTap: onRepost,
                      ),
                      _ActionButton(
                        icon: post.isLiked ? Icons.favorite : Icons.favorite_border,
                        count: post.likeCount,
                        isActive: post.isLiked,
                        activeColor: QubeTheme.danger,
                        onTap: onLike,
                      ),
                      _ActionButton(
                        icon: post.isBookmarked ? Icons.bookmark : Icons.bookmark_border,
                        isActive: post.isBookmarked,
                        activeColor: QubeTheme.primary,
                        onTap: onBookmark,
                      ),
                      _ActionButton(
                        icon: Icons.share_outlined,
                        onTap: () {},
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _ActionButton extends StatelessWidget {
  final IconData icon;
  final int? count;
  final bool isActive;
  final Color? activeColor;
  final VoidCallback? onTap;

  const _ActionButton({
    required this.icon,
    this.count,
    this.isActive = false,
    this.activeColor,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final color = isActive ? (activeColor ?? QubeTheme.primary) : QubeTheme.textSecondary;
    return GestureDetector(
      onTap: onTap,
      child: Row(
        children: [
          Icon(icon, size: 18, color: color),
          if (count != null && count! > 0) ...[
            const SizedBox(width: 4),
            Text(
              _formatCount(count!),
              style: TextStyle(color: color, fontSize: 13),
            ),
          ],
        ],
      ),
    );
  }

  String _formatCount(int count) {
    if (count >= 1000000) return '${(count / 1000000).toStringAsFixed(1)}M';
    if (count >= 1000) return '${(count / 1000).toStringAsFixed(1)}K';
    return count.toString();
  }
}
