var path    = require('path')
var webpack = require('webpack')

module.exports = {
  entry: './src/index',
  output: {
    path: __dirname,
    filename: 'bundle.js'
  },
  resolve: {
    modules: [
        path.resolve('./src'),
        "node_modules"
    ],
    extensions: ['.js', '.jsx']
  },
  module: {
    rules: [
      {
        test: /\.jsx?$/,
        loader: 'babel-loader',
        exclude: /node_modules/,
        query: {
          presets: ['es2015', 'react', 'stage-2']
        }
      }
    ]
  }
}