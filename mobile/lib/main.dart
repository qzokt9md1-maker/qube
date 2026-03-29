import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'config/theme.dart';
import 'config/constants.dart';
import 'providers/auth_provider.dart';
import 'screens/auth/login_screen.dart';
import 'screens/auth/register_screen.dart';
import 'screens/home/main_shell.dart';

void main() {
  runApp(const ProviderScope(child: QubeApp()));
}

class QubeApp extends ConsumerStatefulWidget {
  const QubeApp({super.key});

  @override
  ConsumerState<QubeApp> createState() => _QubeAppState();
}

class _QubeAppState extends ConsumerState<QubeApp> {
  @override
  void initState() {
    super.initState();
    ref.read(authProvider.notifier).init();
  }

  @override
  Widget build(BuildContext context) {
    final auth = ref.watch(authProvider);

    return MaterialApp(
      title: QubeConstants.appName,
      theme: QubeTheme.darkTheme,
      debugShowCheckedModeBanner: false,
      home: auth.isLoggedIn ? const MainShell() : const WelcomeScreen(),
    );
  }
}

class WelcomeScreen extends StatelessWidget {
  const WelcomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: Center(
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 32),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                RichText(
                  text: const TextSpan(
                    children: [
                      TextSpan(
                        text: 'Q',
                        style: TextStyle(fontSize: 64, fontWeight: FontWeight.bold, color: QubeTheme.primary),
                      ),
                      TextSpan(
                        text: 'ube',
                        style: TextStyle(fontSize: 64, fontWeight: FontWeight.bold, color: QubeTheme.textPrimary),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 12),
                const Text(
                  'A social network where you\nnever miss a post.',
                  textAlign: TextAlign.center,
                  style: TextStyle(fontSize: 16, color: QubeTheme.textSecondary),
                ),
                const SizedBox(height: 48),
                SizedBox(
                  width: double.infinity,
                  child: ElevatedButton(
                    onPressed: () => Navigator.push(context, MaterialPageRoute(builder: (_) => const RegisterScreen())),
                    child: const Text('Sign Up', style: TextStyle(fontSize: 16, fontWeight: FontWeight.w600)),
                  ),
                ),
                const SizedBox(height: 12),
                SizedBox(
                  width: double.infinity,
                  child: OutlinedButton(
                    onPressed: () => Navigator.push(context, MaterialPageRoute(builder: (_) => const LoginScreen())),
                    style: OutlinedButton.styleFrom(
                      side: const BorderSide(color: QubeTheme.border),
                      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(24)),
                      padding: const EdgeInsets.symmetric(vertical: 14),
                    ),
                    child: const Text('Log In', style: TextStyle(fontSize: 16, fontWeight: FontWeight.w600, color: QubeTheme.textPrimary)),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
