/**
 * Returns a debounced version of a function
 * @param fn Function to debounce
 * @param delay Delay in milliseconds
 */
export function useDebounce() {
  let timeout: number | null = null;

  // Debounced function
  const debounced = (fn: () => void, delay: number) => {
    if (timeout) {
      clearTimeout(timeout);
    }

    timeout = setTimeout(() => {
      fn();
      timeout = null;
    }, delay);
  };

  // Cancel any pending debounce
  const cancel = () => {
    if (timeout) {
      clearTimeout(timeout);
      timeout = null;
    }
  };

  return { debounced, cancel };
}
