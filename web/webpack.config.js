const path = require('path');
const webpack = require('webpack');
const dotenv = require("dotenv")

dotenv.config()

module.exports = {
  entry: './src/script-index.js',
  output: {
    filename: 'bundle.js',
    path: path.resolve(__dirname, 'static-files'),
    libraryTarget: 'global',
    library: 'script'
  },
  resolve: {
    modules: ['node_modules']
  },
  plugins:[
    new webpack.DefinePlugin({
      // list config variables here!
      CONFIG_API_ASSETS_URL: JSON.stringify(process.env.CONFIG_API_ASSETS_URL),
      CONFIG_API_TW_URL: JSON.stringify(process.env.CONFIG_API_TW_URL),
    })
  ],
  mode: "development",
  module: {
    rules: [
      {
        test: /\.css$/i,
        use: ['style-loader', 'css-loader'],
      },
    ],
  },
};
