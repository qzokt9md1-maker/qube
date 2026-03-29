import 'dart:async';
import 'package:flutter/material.dart';
import '../../config/theme.dart';
import '../../models/user.dart';
import '../../services/api_service.dart';
import '../profile/profile_screen.dart';

class SearchScreen extends StatefulWidget {
  const SearchScreen({super.key});

  @override
  State<SearchScreen> createState() => _SearchScreenState();
}

class _SearchScreenState extends State<SearchScreen> {
  final _api = ApiService();
  final _controller = TextEditingController();
  Timer? _debounce;
  List<QubeUser> _results = [];
  bool _isSearching = false;

  @override
  void dispose() {
    _controller.dispose();
    _debounce?.cancel();
    super.dispose();
  }

  void _onSearchChanged(String query) {
    _debounce?.cancel();
    _debounce = Timer(const Duration(milliseconds: 300), () {
      if (query.trim().isEmpty) {
        setState(() => _results = []);
        return;
      }
      _search(query.trim());
    });
  }

  Future<void> _search(String query) async {
    setState(() => _isSearching = true);
    try {
      final data = await _api.query('SearchUsers', variables: {'query': query, 'limit': 20});
      setState(() {
        _results = (data['searchUsers']['users'] as List)
            .map((u) => QubeUser.fromJson(u))
            .toList();
        _isSearching = false;
      });
    } catch (e) {
      setState(() => _isSearching = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: TextField(
          controller: _controller,
          autofocus: true,
          decoration: InputDecoration(
            hintText: 'Search users...',
            border: InputBorder.none,
            hintStyle: TextStyle(color: QubeTheme.textSecondary.withValues(alpha: 0.6)),
          ),
          style: const TextStyle(fontSize: 16),
          onChanged: _onSearchChanged,
        ),
        actions: [
          if (_controller.text.isNotEmpty)
            IconButton(
              icon: const Icon(Icons.clear),
              onPressed: () {
                _controller.clear();
                setState(() => _results = []);
              },
            ),
        ],
      ),
      body: _isSearching
          ? const Center(child: CircularProgressIndicator())
          : _results.isEmpty
              ? Center(
                  child: Text(
                    _controller.text.isEmpty ? 'Search for users' : 'No results found',
                    style: const TextStyle(color: QubeTheme.textSecondary),
                  ),
                )
              : ListView.builder(
                  itemCount: _results.length,
                  itemBuilder: (context, index) {
                    final user = _results[index];
                    return ListTile(
                      leading: CircleAvatar(
                        backgroundColor: QubeTheme.surface,
                        backgroundImage: user.avatarUrl.isNotEmpty ? NetworkImage(user.avatarUrl) : null,
                        child: user.avatarUrl.isEmpty
                            ? Text(user.displayName[0].toUpperCase())
                            : null,
                      ),
                      title: Row(
                        children: [
                          Text(user.displayName, style: const TextStyle(fontWeight: FontWeight.w600)),
                          if (user.isVerified) ...[
                            const SizedBox(width: 4),
                            const Icon(Icons.verified, size: 16, color: QubeTheme.primary),
                          ],
                        ],
                      ),
                      subtitle: Text('@${user.username}', style: const TextStyle(color: QubeTheme.textSecondary)),
                      trailing: Text(
                        '${user.followerCount} followers',
                        style: const TextStyle(color: QubeTheme.textSecondary, fontSize: 12),
                      ),
                      onTap: () {
                        Navigator.push(
                          context,
                          MaterialPageRoute(builder: (_) => ProfileScreen(username: user.username)),
                        );
                      },
                    );
                  },
                ),
    );
  }
}
