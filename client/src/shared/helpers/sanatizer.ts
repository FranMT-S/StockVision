/**
 * Sanitizes input for SQL remove characters that can be used to inject SQL code
 */
export function sanitizeSQL(input: string | null | undefined): string {
  if (!input) return '';
  return input
  .replace(/'/g, "''")
  .replace(/;/g, '')
  .replace(/\*/g, '')
  .replace(/\//g, '')
  .replace(/\\/g, '')
  .replace(/-/g, ''); 
}

/**
 * Sanitizes input for use in HTML/JS
 */
export function sanitizeJS(input: string | null | undefined): string {
  if (!input) return '';
  return input
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')
    .replace(/\//g, '&#x2F;');
}
