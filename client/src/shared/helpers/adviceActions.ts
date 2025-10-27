import vuetify from "@/plugins/vuetify"

export const getAdviceActionIcon = (action: string) => {
  switch (action) {
    case 'BUY':
      return 'mdi-cart-plus'
    case 'SELL':
      return 'mdi-trending-down'
    case 'HOLD':
      return 'mdi-hand-back-right'
    default:
      return ''
  }
}

export const getAdviceActionColor = (action: string) => {
  switch (action) {
    case 'BUY':
      return {class:'tw-text-green-500', vuetify:'green' }
    case 'SELL':
      return {class:'tw-text-red-500', vuetify:'red' }
    case 'HOLD':
      return {class:'tw-text-gray-500', vuetify:'gray' }
    default:
      return {class:'', vuetify:'' }
  }
}
