import 'package:flutter/material.dart';
import '../../config/theme.dart';
import '../../models/user.dart';
import '../../models/post.dart';
import '../../services/api_service.dart';
import '../../widgets/post_card.dart';

class ProfileScreen extends StatefulWidget {
  final String username;

  const ProfileScreen({super.key, required this.username});

  @override
  State<ProfileScreen> createState() => _ProfileScreenState();
}

class _ProfileScreenState extends State<ProfileScreen> with SingleTickerProviderStateMixin {
  final _api = ApiService();
  late TabController _tabController;
  QubeUser? _user;
  List<QubePost> _posts = [];
  final List<QubePost> _likes = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
    _loadProfile();
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  Future<void> _loadProfile() async {
    try {
      final userData = await _api.query('User', variables: {'username': widget.username});
      final postsData = await _api.query('UserPosts', variables: {'username': widget.username, 'limit': 20});

      setState(() {
        _user = QubeUser.fromJson(userData['user']);
        _posts = (postsData['userPosts']['posts'] as List).map((p) => QubePost.fromJson(p)).toList();
        _isLoading = false;
      });
    } catch (e) {
      setState(() => _isLoading = false);
    }
  }

  Future<void> _toggleFollow() async {
    // TODO: Check if following, then follow/unfollow
    await _api.query('Follow', variables: {'userId': _user!.id});
    _loadProfile();
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading) {
      return const Scaffold(body: Center(child: CircularProgressIndicator()));
    }

    if (_user == null) {
      return Scaffold(
        appBar: AppBar(),
        body: const Center(child: Text('User not found')),
      );
    }

    return Scaffold(
      body: NestedScrollView(
        headerSliverBuilder: (context, _) => [
          SliverAppBar(
            expandedHeight: 200,
            pinned: true,
            flexibleSpace: FlexibleSpaceBar(
              background: _user!.headerUrl.isNotEmpty
                  ? Image.network(_user!.headerUrl, fit: BoxFit.cover)
                  : Container(color: QubeTheme.primary.withValues(alpha: 0.3)),
            ),
          ),
          SliverToBoxAdapter(
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Avatar + Follow button row
                  Transform.translate(
                    offset: const Offset(0, -40),
                    child: Row(
                      crossAxisAlignment: CrossAxisAlignment.end,
                      children: [
                        CircleAvatar(
                          radius: 40,
                          backgroundColor: QubeTheme.background,
                          child: CircleAvatar(
                            radius: 37,
                            backgroundColor: QubeTheme.surface,
                            backgroundImage: _user!.avatarUrl.isNotEmpty ? NetworkImage(_user!.avatarUrl) : null,
                            child: _user!.avatarUrl.isEmpty
                                ? Text(_user!.displayName[0].toUpperCase(), style: const TextStyle(fontSize: 32))
                                : null,
                          ),
                        ),
                        const Spacer(),
                        OutlinedButton(
                          onPressed: _toggleFollow,
                          style: OutlinedButton.styleFrom(
                            side: const BorderSide(color: QubeTheme.border),
                            shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(20)),
                          ),
                          child: const Text('Follow', style: TextStyle(color: QubeTheme.textPrimary)),
                        ),
                      ],
                    ),
                  ),
                  Transform.translate(
                    offset: const Offset(0, -24),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Row(
                          children: [
                            Text(_user!.displayName, style: const TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
                            if (_user!.isVerified) ...[
                              const SizedBox(width: 4),
                              const Icon(Icons.verified, size: 20, color: QubeTheme.primary),
                            ],
                          ],
                        ),
                        Text('@${_user!.username}', style: const TextStyle(color: QubeTheme.textSecondary)),
                        if (_user!.bio.isNotEmpty) ...[
                          const SizedBox(height: 8),
                          Text(_user!.bio, style: const TextStyle(fontSize: 15)),
                        ],
                        const SizedBox(height: 12),
                        Row(
                          children: [
                            _StatItem(count: _user!.followingCount, label: 'Following'),
                            const SizedBox(width: 16),
                            _StatItem(count: _user!.followerCount, label: 'Followers'),
                          ],
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ),
          SliverPersistentHeader(
            pinned: true,
            delegate: _TabBarDelegate(
              TabBar(
                controller: _tabController,
                indicatorColor: QubeTheme.primary,
                labelColor: QubeTheme.textPrimary,
                unselectedLabelColor: QubeTheme.textSecondary,
                tabs: const [
                  Tab(text: 'Posts'),
                  Tab(text: 'Replies'),
                  Tab(text: 'Likes'),
                ],
              ),
            ),
          ),
        ],
        body: TabBarView(
          controller: _tabController,
          children: [
            // Posts
            ListView.builder(
              itemCount: _posts.length,
              itemBuilder: (context, index) => PostCard(post: _posts[index]),
            ),
            // Replies (placeholder)
            const Center(child: Text('Replies', style: TextStyle(color: QubeTheme.textSecondary))),
            // Likes
            ListView.builder(
              itemCount: _likes.length,
              itemBuilder: (context, index) => PostCard(post: _likes[index]),
            ),
          ],
        ),
      ),
    );
  }
}

class _StatItem extends StatelessWidget {
  final int count;
  final String label;
  const _StatItem({required this.count, required this.label});

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Text('$count', style: const TextStyle(fontWeight: FontWeight.bold)),
        const SizedBox(width: 4),
        Text(label, style: const TextStyle(color: QubeTheme.textSecondary)),
      ],
    );
  }
}

class _TabBarDelegate extends SliverPersistentHeaderDelegate {
  final TabBar tabBar;
  _TabBarDelegate(this.tabBar);

  @override
  Widget build(BuildContext context, double shrinkOffset, bool overlapsContent) {
    return Container(color: QubeTheme.background, child: tabBar);
  }

  @override
  double get maxExtent => tabBar.preferredSize.height;
  @override
  double get minExtent => tabBar.preferredSize.height;
  @override
  bool shouldRebuild(covariant _TabBarDelegate oldDelegate) => false;
}
