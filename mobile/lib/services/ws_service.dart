import 'dart:async';
import 'dart:convert';
import 'package:web_socket_channel/web_socket_channel.dart';
import '../config/constants.dart';

class WsService {
  static final WsService _instance = WsService._internal();
  factory WsService() => _instance;
  WsService._internal();

  WebSocketChannel? _channel;
  final _controller = StreamController<Map<String, dynamic>>.broadcast();

  Stream<Map<String, dynamic>> get stream => _controller.stream;

  void connect(String token) {
    _channel = WebSocketChannel.connect(
      Uri.parse('${QubeConstants.wsUrl}?token=$token'),
    );
    _channel!.stream.listen(
      (data) {
        try {
          final event = jsonDecode(data as String) as Map<String, dynamic>;
          _controller.add(event);
        } catch (_) {}
      },
      onDone: () {
        // Reconnect after 3 seconds
        Future.delayed(const Duration(seconds: 3), () => connect(token));
      },
    );
  }

  void send(Map<String, dynamic> event) {
    _channel?.sink.add(jsonEncode(event));
  }

  void sendTyping(String conversationId) {
    send({
      'type': 'typing',
      'payload': {'conversation_id': conversationId},
    });
  }

  void disconnect() {
    _channel?.sink.close();
    _channel = null;
  }
}
