/**
 * Get today's date as a YYYY-MM-DD string in the local timezone.
 *
 * `new Date().toISOString().split("T")[0]` returns UTC date, which
 * is wrong in JST (UTC+9) between 00:00-08:59 JST.
 * `toLocaleDateString("sv-SE")` returns YYYY-MM-DD in local time.
 */
export function getLocalDateString(): string {
	return new Date().toLocaleDateString("sv-SE");
}
