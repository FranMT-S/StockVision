import { describe, it, expect } from 'vitest'
import { sanitizeSQL, sanitizeJS } from '@/shared/helpers/sanatizer'

describe('sanitizer helper', () => {
  describe('sanitizeSQL', () => {
    it('returns empty string for null or undefined input', () => {
      expect(sanitizeSQL(null)).toBe('')
      expect(sanitizeSQL(undefined)).toBe('')
    })



    it('removes semicolons', () => {
      expect(sanitizeSQL("SELECT * FROM users;")).toBe("SELECT  FROM users")
      expect(sanitizeSQL("DROP TABLE users;")).toBe("DROP TABLE users")
    })

    it('removes SQL comment characters', () => {
      expect(sanitizeSQL("test--comment")).toBe("testcomment")
      expect(sanitizeSQL("test/*comment*/")).toBe("testcomment")
    })

  })

  describe('sanitizeJS', () => {
    it('returns empty string for null or undefined input', () => {
      expect(sanitizeJS(null)).toBe('')
      expect(sanitizeJS(undefined)).toBe('')
    })

    it('escapes HTML special characters', () => {
      expect(sanitizeJS('<script>alert(1)</script>'))
        .toBe('&lt;script&gt;alert(1)&lt;&#x2F;script&gt;')
      
      expect(sanitizeJS('"quoted" & \'string\''))
        .toBe('&quot;quoted&quot; &amp; &#39;string&#39;')
    })

    it('escapes forward slashes', () => {
      expect(sanitizeJS('</script>')).toBe('&lt;&#x2F;script&gt;')
      expect(sanitizeJS('https://example.com')).toBe('https:&#x2F;&#x2F;example.com')
    })

    it('handles complex XSS attempts', () => {
      const xss = '"><img src=x onerror=alert(1)>'
      const safe = '&quot;&gt;&lt;img src=x onerror=alert(1)&gt;'
      expect(sanitizeJS(xss)).toBe(safe)
    })
  })
})
