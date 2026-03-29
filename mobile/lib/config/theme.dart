import 'package:flutter/material.dart';

class QubeTheme {
  static const primary = Color(0xFFC9A96E);
  static const primaryDark = Color(0xFFB8955D);
  static const background = Color(0xFF000000);
  static const surface = Color(0xFF16181C);
  static const surfaceHover = Color(0xFF1D1F23);
  static const border = Color(0xFF2F3336);
  static const textPrimary = Color(0xFFE7E9EA);
  static const textSecondary = Color(0xFF71767B);
  static const danger = Color(0xFFF91880);
  static const success = Color(0xFF00BA7C);
  static const like = Color(0xFFF91880);
  static const repost = Color(0xFF00BA7C);

  static ThemeData get darkTheme => ThemeData(
        brightness: Brightness.dark,
        scaffoldBackgroundColor: background,
        primaryColor: primary,
        colorScheme: const ColorScheme.dark(
          primary: primary,
          secondary: primaryDark,
          surface: surface,
          error: danger,
        ),
        appBarTheme: const AppBarTheme(
          backgroundColor: background,
          elevation: 0,
          centerTitle: false,
          titleTextStyle: TextStyle(
            color: textPrimary,
            fontSize: 20,
            fontWeight: FontWeight.bold,
          ),
        ),
        dividerTheme: const DividerThemeData(
          color: border,
          thickness: 0.5,
        ),
        inputDecorationTheme: InputDecorationTheme(
          filled: true,
          fillColor: surface,
          border: OutlineInputBorder(
            borderRadius: BorderRadius.circular(12),
            borderSide: const BorderSide(color: border),
          ),
          enabledBorder: OutlineInputBorder(
            borderRadius: BorderRadius.circular(12),
            borderSide: const BorderSide(color: border),
          ),
          focusedBorder: OutlineInputBorder(
            borderRadius: BorderRadius.circular(12),
            borderSide: const BorderSide(color: primary),
          ),
        ),
        elevatedButtonTheme: ElevatedButtonThemeData(
          style: ElevatedButton.styleFrom(
            backgroundColor: primary,
            foregroundColor: Colors.black,
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(24),
            ),
            padding: const EdgeInsets.symmetric(horizontal: 32, vertical: 14),
          ),
        ),
      );
}
