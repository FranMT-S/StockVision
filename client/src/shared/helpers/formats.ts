
export const humanizeNumberFormat = (num: number, decimals = 1, allowTrillions = false): string => {
  if (allowTrillions && num >= 1e12) {
    return (num / 1e12).toFixed(decimals) + 'T'
  } else if (num >= 1e9) {
    return (num / 1e9).toFixed(decimals) + 'B'
  } else if (num >= 1e6) {
    return (num / 1e6).toFixed(decimals) + 'M'
  } else if (num >= 1e3) {
    return (num / 1e3).toFixed(decimals) + 'K'
  }
  return num.toString()
}