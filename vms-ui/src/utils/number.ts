const prefixes = [
  {value: 1e24, symbol: 'Y'},  // yotta
  {value: 1e21, symbol: 'Z'},  // zetta
  {value: 1e18, symbol: 'E'},  // exa
  {value: 1e15, symbol: 'P'},  // peta
  {value: 1e12, symbol: 'T'},  // tera
  {value: 1e9, symbol: 'G'},   // giga
  {value: 1e6, symbol: 'M'},   // mega
  {value: 1e3, symbol: 'k'},   // kilo
  {value: 1, symbol: ''},      // base
  {value: 1e-3, symbol: 'm'},  // milli
  {value: 1e-6, symbol: 'µ'},  // micro
  {value: 1e-9, symbol: 'n'},  // nano
  {value: 1e-12, symbol: 'p'}, // pico
  {value: 1e-15, symbol: 'f'}, // femto
  {value: 1e-18, symbol: 'a'}, // atto
  {value: 1e-21, symbol: 'z'}, // zepto
  {value: 1e-24, symbol: 'y'}  // yocto
];

export function formatSIUnit(value: any, baseUnit = '', decimals = 2) {
  if (!value || value === 0) return `0 ${baseUnit}`;

  const absValue = Math.abs(value);

  for (let i = 0; i < prefixes.length; i++) {
    if (absValue >= prefixes[i].value) {
      const scaled = value / prefixes[i].value;
      return `${scaled.toFixed(decimals)} ${prefixes[i].symbol}${baseUnit}`;
    }
  }

  return `${value.toExponential(decimals)} ${baseUnit}`;
}
