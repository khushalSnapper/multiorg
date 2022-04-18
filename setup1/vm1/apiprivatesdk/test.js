// function that takes a function and returns a function
const memoize = (func) => {
    // a cache of results
    const results = {};
    // return a function for the cache of results
    return (...args) => {
      // a JSON key to save the results cache
      const argsKey = JSON.stringify(args);
      // execute `func` only if there is no cached value of clumsysquare()
      if (!results[argsKey]) {
        // store the return value of clumsysquare()
        results[argsKey] = func(...args);
      }
      // return the cached results
      return results[argsKey];
    };
  };
  
  // wrap clumsysquare() in memoize()
  const clumsysquare = memoize(num => {
      let result = 0;
      for (let i = 1; i <= num; i++) {
          for (let j = 1; j <= num; j++) {
              result++;
          }
      }
      return result;
  });
  
  console.log(clumsysquare(190));
  console.log(clumsysquare(799));
  console.log(clumsysquare(4000));
  console.log(clumsysquare(7467));
  console.log(clumsysquare(9666));