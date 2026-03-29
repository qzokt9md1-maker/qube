import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/user.dart';
import '../services/api_service.dart';

class AuthState {
  final QubeUser? user;
  final bool isLoading;
  final String? error;

  AuthState({this.user, this.isLoading = false, this.error});

  AuthState copyWith({QubeUser? user, bool? isLoading, String? error}) {
    return AuthState(
      user: user ?? this.user,
      isLoading: isLoading ?? this.isLoading,
      error: error,
    );
  }

  bool get isLoggedIn => user != null;
}

class AuthNotifier extends StateNotifier<AuthState> {
  final ApiService _api = ApiService();

  AuthNotifier() : super(AuthState());

  Future<void> init() async {
    await _api.init();
    if (_api.isLoggedIn) {
      try {
        final data = await _api.query('Me');
        state = AuthState(user: QubeUser.fromJson(data['me']));
      } catch (_) {
        await _api.clearTokens();
        state = AuthState();
      }
    }
  }

  Future<void> register(String username, String displayName, String email, String password) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final data = await _api.query('Register', variables: {
        'input': {
          'username': username,
          'displayName': displayName,
          'email': email,
          'password': password,
        },
      });
      final result = data['register'];
      await _api.setTokens(result['accessToken'], result['refreshToken']);
      state = AuthState(user: QubeUser.fromJson(result['user']));
    } on ApiException catch (e) {
      state = state.copyWith(isLoading: false, error: e.message);
    }
  }

  Future<void> login(String email, String password) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final data = await _api.query('Login', variables: {
        'input': {'email': email, 'password': password},
      });
      final result = data['login'];
      await _api.setTokens(result['accessToken'], result['refreshToken']);
      state = AuthState(user: QubeUser.fromJson(result['user']));
    } on ApiException catch (e) {
      state = state.copyWith(isLoading: false, error: e.message);
    }
  }

  Future<void> logout() async {
    await _api.clearTokens();
    state = AuthState();
  }
}

final authProvider = StateNotifierProvider<AuthNotifier, AuthState>((ref) {
  return AuthNotifier();
});
