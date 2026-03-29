import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../config/theme.dart';
import '../../models/post.dart';
import '../../services/api_service.dart';
import '../../widgets/post_card.dart';

class TimelineScreen extends ConsumerStatefulWidget {
  const TimelineScreen({super.key});

  @override
  ConsumerState<TimelineScreen> createState() => _TimelineScreenState();
}

class _TimelineScreenState extends ConsumerState<TimelineScreen> {
  final _api = ApiService();
  final _scrollController = ScrollController();
  List<QubePost> _posts = [];
  bool _isLoading = false;
  bool _hasMore = true;
  String? _cursor;
  int _unreadCount = 0;

  @override
  void initState() {
    super.initState();
    _loadTimeline();
    _scrollController.addListener(_onScroll);
  }

  @override
  void dispose() {
    _scrollController.dispose();
    super.dispose();
  }

  void _onScroll() {
    if (_scrollController.position.pixels >= _scrollController.position.maxScrollExtent - 200) {
      _loadMore();
    }
  }

  Future<void> _loadTimeline() async {
    setState(() => _isLoading = true);
    try {
      final data = await _api.query('Timeline', variables: {'limit': 20});
      final timeline = data['timeline'];
      setState(() {
        _posts = (timeline['posts'] as List).map((p) => QubePost.fromJson(p)).toList();
        _hasMore = timeline['hasMore'] as bool;
        _cursor = timeline['cursor'] as String?;
        _unreadCount = timeline['unreadCount'] as int? ?? 0;
        _isLoading = false;
      });
      // Mark as read
      if (_posts.isNotEmpty) {
        _api.query('UpdateTimelineCursor', variables: {'lastSeenPostId': _posts.first.id});
      }
    } catch (e) {
      setState(() => _isLoading = false);
    }
  }

  Future<void> _loadMore() async {
    if (_isLoading || !_hasMore) return;
    setState(() => _isLoading = true);
    try {
      final data = await _api.query('Timeline', variables: {'limit': 20, 'cursor': _cursor});
      final timeline = data['timeline'];
      final newPosts = (timeline['posts'] as List).map((p) => QubePost.fromJson(p)).toList();
      setState(() {
        _posts.addAll(newPosts);
        _hasMore = timeline['hasMore'] as bool;
        _cursor = timeline['cursor'] as String?;
        _isLoading = false;
      });
    } catch (e) {
      setState(() => _isLoading = false);
    }
  }

  Future<void> _onLike(QubePost post) async {
    final op = post.isLiked ? 'UnlikePost' : 'LikePost';
    await _api.query(op, variables: {'postId': post.id});
    _loadTimeline();
  }

  Future<void> _onRepost(QubePost post) async {
    await _api.query('Repost', variables: {'postId': post.id});
    _loadTimeline();
  }

  Future<void> _onBookmark(QubePost post) async {
    final op = post.isBookmarked ? 'UnbookmarkPost' : 'BookmarkPost';
    await _api.query(op, variables: {'postId': post.id});
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: RichText(
          text: const TextSpan(children: [
            TextSpan(text: 'Q', style: TextStyle(color: QubeTheme.primary, fontSize: 24, fontWeight: FontWeight.bold)),
            TextSpan(text: 'ube', style: TextStyle(color: QubeTheme.textPrimary, fontSize: 24, fontWeight: FontWeight.bold)),
          ]),
        ),
        actions: [
          if (_unreadCount > 0)
            Padding(
              padding: const EdgeInsets.only(right: 16),
              child: Center(
                child: Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                  decoration: BoxDecoration(
                    color: QubeTheme.primary,
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Text('$_unreadCount new', style: const TextStyle(fontSize: 12)),
                ),
              ),
            ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: _loadTimeline,
        child: _posts.isEmpty && !_isLoading
            ? const Center(
                child: Text('Follow someone to see their posts!', style: TextStyle(color: QubeTheme.textSecondary)),
              )
            : ListView.builder(
                controller: _scrollController,
                itemCount: _posts.length + (_isLoading ? 1 : 0),
                itemBuilder: (context, index) {
                  if (index == _posts.length) {
                    return const Padding(
                      padding: EdgeInsets.all(16),
                      child: Center(child: CircularProgressIndicator()),
                    );
                  }
                  final post = _posts[index];
                  return PostCard(
                    post: post,
                    onLike: () => _onLike(post),
                    onRepost: () => _onRepost(post),
                    onBookmark: () => _onBookmark(post),
                    onTap: () {
                      Navigator.push(context, MaterialPageRoute(
                        builder: (_) => PostDetailScreen(postId: post.id),
                      ));
                    },
                    onUserTap: () {
                      Navigator.push(context, MaterialPageRoute(
                        builder: (_) => ProfileScreenNav(username: post.user.username),
                      ));
                    },
                  );
                },
              ),
      ),
    );
  }
}

// Placeholder navigations - will be replaced by proper routing
class PostDetailScreen extends StatelessWidget {
  final String postId;
  const PostDetailScreen({super.key, required this.postId});
  @override
  Widget build(BuildContext context) => Scaffold(appBar: AppBar(title: const Text('Post')));
}

class ProfileScreenNav extends StatelessWidget {
  final String username;
  const ProfileScreenNav({super.key, required this.username});
  @override
  Widget build(BuildContext context) => Scaffold(appBar: AppBar(title: Text('@$username')));
}
