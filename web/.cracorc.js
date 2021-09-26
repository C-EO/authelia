const isCoverage = process.env.COVERAGE === "true";
const babelPlugins = isCoverage ? ["babel-plugin-istanbul"] : [];

module.exports = {
    babel: {
        plugins: babelPlugins,
    },
    plugins: [
        {
            plugin: require("craco-alias"),
            options: {
                source: "tsconfig",
                baseUrl: "./src",
                tsConfigPath: "./tsconfig.aliases.json",
            },
        },
        {
            plugin: require("craco-esbuild"),
            options: {
                enableSvgr: true,
                esbuildMinimizerOptions: {
                    target: "es6",
                    css: true,
                    minify: true,
                },
                skipEsbuildJest: true,
            },
        },
    ],
};
