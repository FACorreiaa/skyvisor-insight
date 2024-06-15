module.exports = {
  plugins: [
    require('postcss-import'),
    require('postcss-url')({
      url: 'copy',
      useHash: true,
      assetsPath: '../fonts',
    }),
    require('tailwindcss'),
    require('autoprefixer')
  ],
};
