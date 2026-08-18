import 'package:flutter_test/flutter_test.dart';

import 'package:kulladhisab_app/main.dart';

void main() {
  test('summarize blends cost per cup and margin', () {
    final s = summarize([DayLog(600, 100, 10), DayLog(800, 100, 10)]);
    expect(s.avgCost, closeTo(7, 1e-9));
    expect(s.margin, closeTo(3, 1e-9));
  });

  testWidgets('shows the cost-per-cup header', (tester) async {
    await tester.pumpWidget(const KulladhisabApp());
    expect(find.text('Average cost per cup'), findsOneWidget);
  });
}
