/* eslint strict: 0 */
const webpack = require("webpack");
const ExtractTextPlugin = require("extract-text-webpack-plugin");
const path = require("path");

const buildQS = function(props) {
	"use strict";
	let out = [];
	for (let prop in props) {
		out.push(props[prop].map((item) => `${prop}[]=${item}`).join(","));
	}
	return out.join(",");
}

const babelcmd = "babel?" + buildQS({
	presets: [
		"es2015",
		"stage-0"
	],
	plugins: [
		"transform-decorators-legacy"
	]
});

module.exports = {
	babelcmd: babelcmd,
	devtool: "sourceMap",
	entry: {
		app: ["./static-src/index"]
	},
	module: {
		loaders: [{
			test: /\.js?$/,
			loaders: [babelcmd],
			exclude: /node_modules/
		},{
			test: /\.eot|.otf|.woff|\.ttf/,
			loader: "file"
		},{
			test: /\.scss$/,
			loader: ExtractTextPlugin.extract(
				"style",
				"css",
				"sass?sourceMap"
			)
		},{
			test: /\.css$/,
			loader: ExtractTextPlugin.extract(
				"style",
				"css"
			)
		}]
	},
	output: {
		path: path.join(__dirname, "static"),
		publicPath: "/static",
		filename: "bundle.js"
	},
	resolve: {
		extensions: ["", ".js", ".scss", ".css"],
		root: [path.join(__dirname, "./static-src")]
	},
	plugins: [
		new webpack.optimize.OccurenceOrderPlugin(),
		new webpack.DefinePlugin({
			"__DEV__": false,
			"process.env": {
				"NODE_ENV": JSON.stringify("production")
			}
		}),
		new webpack.optimize.UglifyJsPlugin({
			compressor: {
				screw_ie8: true,
				warnings: false
			}
		}),
		new ExtractTextPlugin("style.css", { allChunks: true })
	],
	externals: []
};