import 'package:flutter/material.dart';
import '../../config/theme.dart';
import '../post/compose_screen.dart';
import '../search/search_screen.dart';
import '../notifications/notifications_screen.dart';
import '../dm/conversations_screen.dart';
import 'timeline_screen.dart';

class MainShell extends StatefulWidget {
  const MainShell({super.key});

  @override
  State<MainShell> createState() => _MainShellState();
}

class _MainShellState extends State<MainShell> {
  int _currentIndex = 0;

  final _screens = const [
    TimelineScreen(),
    SearchScreen(),
    SizedBox(), // Compose placeholder
    NotificationsScreen(),
    ConversationsScreen(),
  ];

  void _onTap(int index) {
    if (index == 2) {
      // Compose
      Navigator.push(
        context,
        MaterialPageRoute(
          fullscreenDialog: true,
          builder: (_) => const ComposeScreen(),
        ),
      ).then((posted) {
        if (posted == true) {
          setState(() => _currentIndex = 0);
        }
      });
      return;
    }
    setState(() => _currentIndex = index);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: IndexedStack(
        index: _currentIndex > 2 ? _currentIndex - 1 : _currentIndex,
        children: [
          _screens[0],
          _screens[1],
          _screens[3],
          _screens[4],
        ],
      ),
      bottomNavigationBar: Container(
        decoration: const BoxDecoration(
          border: Border(top: BorderSide(color: QubeTheme.border, width: 0.5)),
        ),
        child: BottomNavigationBar(
          currentIndex: _currentIndex,
          onTap: _onTap,
          type: BottomNavigationBarType.fixed,
          backgroundColor: QubeTheme.background,
          selectedItemColor: QubeTheme.primary,
          unselectedItemColor: QubeTheme.textSecondary,
          showSelectedLabels: false,
          showUnselectedLabels: false,
          items: const [
            BottomNavigationBarItem(icon: Icon(Icons.home_outlined), activeIcon: Icon(Icons.home), label: 'Home'),
            BottomNavigationBarItem(icon: Icon(Icons.search), activeIcon: Icon(Icons.search), label: 'Search'),
            BottomNavigationBarItem(icon: Icon(Icons.add_circle_outline, size: 32), label: 'Post'),
            BottomNavigationBarItem(icon: Icon(Icons.notifications_outlined), activeIcon: Icon(Icons.notifications), label: 'Notifications'),
            BottomNavigationBarItem(icon: Icon(Icons.mail_outline), activeIcon: Icon(Icons.mail), label: 'Messages'),
          ],
        ),
      ),
    );
  }
}
