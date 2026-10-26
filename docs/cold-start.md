# Cold start · ideal-weight

`make sam-invoke-iw` (SAM CLI 1.170.0, `provided.al2023`, x86_64, 256 МБ), 2026-10-26.

| Подія | statusCode | Init Duration | Duration |
|---|---|---|---|
| `events/ideal-weight/valid.json` | 200 | 0.13 мс | 17.09 мс |
| `events/ideal-weight/invalid-json.json` | 400 | 0.06 мс | 17.38 мс |
| `events/ideal-weight/missing-field.json` | 400 | 0.04 мс | 7.04 мс |

Init Duration у `sam local` — це час емулятора (RIE), а не справжній холодний старт: запуск
контейнера (~5 с на виклик) відбувається ДО того, як емулятор починає відлік. «Max Memory Used»
теж не вимірюється (показує ліміт). Для порівняння — Floci,
`scripts/bench-cold-start.sh 5 libratrack-ideal-weight` (час виклику через AWS CLI):
cold ≈ 1460 мс, warm ≈ 772 мс (більшість warm — старт самого AWS CLI).
Жодне з чисел не переноситься на AWS (Firecracker) — дивіться на відносну різницю.
