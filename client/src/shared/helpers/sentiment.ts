
export const getSentimentIcon = (sentiment: string): string => {
  switch (sentiment.toLowerCase()) {
    case 'positive':
      return 'mdi-emoticon-happy'
    case 'negative':
      return 'mdi-emoticon-sad'
    case 'neutral':
      return 'mdi-emoticon-neutral'
    default:
      return ''
  }
}


export const getSentimentColor = (sentiment: string): string => {
  switch (sentiment.toLowerCase()) {
    case 'positive':
      return 'green'
    case 'negative':
      return 'red'
    default:
      return 'orange'
  }
}