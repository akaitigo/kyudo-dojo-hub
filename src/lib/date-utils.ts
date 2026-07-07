/**
 * ローカルタイムゾーン基準の日付ユーティリティ。
 *
 * `new Date("2026-04-05")` は ISO の日付のみ文字列を UTC 起点として解釈するため、
 * JST (UTC+9) の 0:00〜8:59 の間は「今日」が UTC 上ではまだ前日となり、
 * 日付比較が 1 日ずれる。これを避けるため、比較は必ずローカル日付文字列
 * (`YYYY-MM-DD`) 同士で行う。`toLocaleDateString("sv-SE")` はロケール "sv-SE"
 * が `YYYY-MM-DD` 形式を返す性質を利用してローカル日付を得る。
 */

/**
 * 指定した日時（省略時は現在時刻）をローカルタイムゾーンの
 * `YYYY-MM-DD` 文字列で返す。
 */
export function getLocalDateString(date: Date = new Date()): string {
	return date.toLocaleDateString("sv-SE");
}

/**
 * 指定した日時（省略時は現在時刻）の 1 年前をローカルタイムゾーンの
 * `YYYY-MM-DD` 文字列で返す。
 */
export function getOneYearAgoDateString(from: Date = new Date()): string {
	const d = new Date(from);
	d.setFullYear(d.getFullYear() - 1);
	return getLocalDateString(d);
}
