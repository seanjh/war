import * as esbuild from 'esbuild'

await esbuild.build({
    entryPoints: ['./client/main.js'],
    bundle: true,
    minify: true,
    outfile: './assets/main.js',
    target: ["chrome58", "firefox57", "safari11", "edge16"],
})
