import { fileURLToPath } from 'node:url'
import { mergeConfig, defineConfig, configDefaults } from 'vitest/config'
import viteConfig from './vite.config'

import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
​
import path from 'path'
// gzip插件
import viteCompression from 'vite-plugin-compression'
// mock插件
import { viteMockServe } from 'vite-plugin-mock'

const resolve = (dir: string) => path.resolve(__dirname, dir)

export default mergeConfig(
  viteConfig,
  defineConfig({
    test: {
      environment: 'jsdom',
      exclude: [...configDefaults.exclude, 'e2e/*'],
      root: fileURLToPath(new URL('./', import.meta.url))
    },
    base: './', //打包路径
    publicDir: resolve('public'), //静态资源服务的文件夹
    plugins: [
        vue(),
        vueJsx(),
        // gzip压缩 生产环境生成 .gz 文件
        viteCompression({
            verbose: true,
            disable: false,
            threshold: 10240,
            algorithm: 'gzip',
            ext: '.gz',
        }),
        //mock
        viteMockServe({
            mockPath: './mocks', // 解析，路径可根据实际变动
            localEnabled: true, // 此处可以手动设置为true，也可以根据官方文档格式
        }),
    ],
    // 配置别名
    resolve: {
        alias: {
            '@': resolve('src'),
        },
        // 导入时想要省略的扩展名列表
        extensions: ['.mjs', '.js', '.ts', '.jsx', '.tsx', '.json', '.vue'],
    },
    css: {
        // css预处理器
        preprocessorOptions: {
            scss: {
                 additionalData: '@import "@/assets/styles/global.scss";@import "@/assets/styles/reset.scss";',
            },
        },
    },
    //启动服务配置
    server: {
        host: '0.0.0.0',
        port: 8081,
        open: true, // 自动在浏览器打开
        proxy: {},
    },
    // 打包配置
    build: {
        //浏览器兼容性  "esnext"|"modules"
        target: 'modules',
        //指定输出路径
        outDir: 'build',
        //生成静态资源的存放路径
        assetsDir: 'assets',
        //启用/禁用 CSS 代码拆分
        cssCodeSplit: true,
        sourcemap: false,
        assetsInlineLimit: 10240,
        // 打包环境移除console.log, debugger
        minify: 'terser',
        terserOptions: {
            compress: {
                drop_console: true,
                drop_debugger: true,
            },
        },
        rollupOptions: {
            input: {
                main: resolve('index.html'),
            },
            output: {
                entryFileNames: `js/[name]-[hash].js`,
                chunkFileNames: `js/[name]-[hash].js`,
                assetFileNames: `[ext]/[name]-[hash].[ext]`,
            },
        },
    },
  })
)
