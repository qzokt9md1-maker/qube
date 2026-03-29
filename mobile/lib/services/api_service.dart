import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import '../config/constants.dart';

class ApiService {
  static final ApiService _instance = ApiService._internal();
  factory ApiService() => _instance;
  ApiService._internal();

  final _storage = const FlutterSecureStorage();
  String? _accessToken;
  String? _refreshToken;

  Future<void> init() async {
    _accessToken = await _storage.read(key: 'access_token');
    _refreshToken = await _storage.read(key: 'refresh_token');
  }

  Future<void> setTokens(String accessToken, String refreshToken) async {
    _accessToken = accessToken;
    _refreshToken = refreshToken;
    await _storage.write(key: 'access_token', value: accessToken);
    await _storage.write(key: 'refresh_token', value: refreshToken);
  }

  Future<void> clearTokens() async {
    _accessToken = null;
    _refreshToken = null;
    await _storage.deleteAll();
  }

  bool get isLoggedIn => _accessToken != null;

  Future<Map<String, dynamic>> query(String operationName, {Map<String, dynamic>? variables}) async {
    final response = await http.post(
      Uri.parse(QubeConstants.apiUrl),
      headers: {
        'Content-Type': 'application/json',
        if (_accessToken != null) 'Authorization': 'Bearer $_accessToken',
      },
      body: jsonEncode({
        'operationName': operationName,
        'query': '',
        'variables': variables ?? {},
      }),
    );

    final data = jsonDecode(response.body) as Map<String, dynamic>;

    if (data['errors'] != null) {
      final errors = data['errors'] as List;
      if (errors.isNotEmpty) {
        final message = errors[0]['message'] as String;
        if (message == 'unauthorized' && _refreshToken != null) {
          final refreshed = await _tryRefresh();
          if (refreshed) {
            return query(operationName, variables: variables);
          }
        }
        throw ApiException(message);
      }
    }

    return data['data'] as Map<String, dynamic>;
  }

  Future<bool> _tryRefresh() async {
    try {
      final response = await http.post(
        Uri.parse(QubeConstants.apiUrl),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({
          'operationName': 'RefreshToken',
          'query': '',
          'variables': {'token': _refreshToken},
        }),
      );
      final data = jsonDecode(response.body) as Map<String, dynamic>;
      if (data['data'] != null && data['data']['refreshToken'] != null) {
        final result = data['data']['refreshToken'];
        await setTokens(result['accessToken'], result['refreshToken']);
        return true;
      }
    } catch (_) {}
    await clearTokens();
    return false;
  }
}

class ApiException implements Exception {
  final String message;
  ApiException(this.message);

  @override
  String toString() => message;
}
