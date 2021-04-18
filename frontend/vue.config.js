module.exports = {
  chainWebpack: config => {
    config
      .plugin("html")
      .tap(args => {
        args[0].title = "Taliesin Millhouse";
        return args;
      });

    // Load Go files as raw text.
    config.module.rule("raw").test(/\.go$/).use("raw-loader").loader("raw-loader").end();
  }
}