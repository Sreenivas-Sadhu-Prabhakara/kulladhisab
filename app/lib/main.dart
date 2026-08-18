import 'package:flutter/material.dart';

void main() => runApp(const KulladhisabApp());

/// Kulladhisab — tea-stall cost per cup. Log daily consumable spend against cups
/// served; the selling price is fixed by habit, so the real cost per cup is what
/// moves. Mirrors the Go journal service.
class KulladhisabApp extends StatelessWidget {
  const KulladhisabApp({super.key});
  @override
  Widget build(BuildContext context) => MaterialApp(
        title: 'Kulladhisab',
        debugShowCheckedModeBanner: false,
        theme: ThemeData(colorSchemeSeed: const Color(0xFF8A5A2B), useMaterial3: true),
        home: const HomePage(),
      );
}

class DayLog {
  final double spend, cups, price;
  DayLog(this.spend, this.cups, this.price);
  double get costPerCup => cups > 0 ? spend / cups : 0;
}

class Summary {
  final double avgCost, avgPrice, margin;
  const Summary(this.avgCost, this.avgPrice, this.margin);
}

/// summarize blends logged days into an overall cost/cup and margin.
Summary summarize(List<DayLog> days) {
  if (days.isEmpty) return const Summary(0, 0, 0);
  double spend = 0, cups = 0, price = 0;
  for (final d in days) {
    spend += d.spend;
    cups += d.cups;
    price += d.price;
  }
  final avgCost = cups > 0 ? spend / cups : 0.0;
  final avgPrice = price / days.length;
  return Summary(avgCost, avgPrice, avgPrice - avgCost);
}

class HomePage extends StatefulWidget {
  const HomePage({super.key});
  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  final _days = <DayLog>[];
  final _spend = TextEditingController();
  final _cups = TextEditingController();
  final _price = TextEditingController(text: '10');

  void _add() {
    final s = double.tryParse(_spend.text.trim()) ?? 0;
    final c = double.tryParse(_cups.text.trim()) ?? 0;
    final p = double.tryParse(_price.text.trim()) ?? 0;
    if (s <= 0 || c <= 0) return;
    setState(() {
      _days.insert(0, DayLog(s, c, p));
      _spend.clear();
      _cups.clear();
    });
  }

  @override
  Widget build(BuildContext context) {
    final sum = summarize(_days);
    String m(num v) => '₹${v.toStringAsFixed(2)}';
    return Scaffold(
      appBar: AppBar(
        title: const Text('Kulladhisab · cost per cup'),
        backgroundColor: Theme.of(context).colorScheme.primaryContainer,
      ),
      body: Column(children: [
        Container(
          width: double.infinity,
          color: Theme.of(context).colorScheme.primaryContainer,
          padding: const EdgeInsets.all(16),
          child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
            const Text('Average cost per cup'),
            Text(m(sum.avgCost), style: const TextStyle(fontSize: 30, fontWeight: FontWeight.bold)),
            Text('Margin ${m(sum.margin)} at ${m(sum.avgPrice)}/cup', style: const TextStyle(fontSize: 13)),
          ]),
        ),
        Padding(
          padding: const EdgeInsets.all(12),
          child: Row(children: [
            Expanded(child: TextField(controller: _spend, keyboardType: TextInputType.number, decoration: const InputDecoration(labelText: 'Day spend ₹', border: OutlineInputBorder()))),
            const SizedBox(width: 8),
            Expanded(child: TextField(controller: _cups, keyboardType: TextInputType.number, decoration: const InputDecoration(labelText: 'Cups sold', border: OutlineInputBorder()))),
            const SizedBox(width: 8),
            SizedBox(width: 80, child: TextField(controller: _price, keyboardType: TextInputType.number, decoration: const InputDecoration(labelText: '₹/cup', border: OutlineInputBorder()))),
            const SizedBox(width: 8),
            FilledButton(onPressed: _add, child: const Text('Log')),
          ]),
        ),
        const Divider(),
        Expanded(
          child: ListView(children: [
            for (final d in _days)
              ListTile(
                dense: true,
                title: Text('${m(d.spend)} over ${d.cups.toStringAsFixed(0)} cups'),
                trailing: Text('${m(d.costPerCup)}/cup', style: const TextStyle(fontWeight: FontWeight.bold)),
              ),
          ]),
        ),
      ]),
    );
  }
}
