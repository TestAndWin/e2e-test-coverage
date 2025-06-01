import pluginVue from 'eslint-plugin-vue'
import { defineConfigWithVueTs, vueTsConfigs } from '@vue/eslint-config-typescript'
import configPrettier from '@vue/eslint-config-prettier/skip-formatting'
import js from '@eslint/js'

export default defineConfigWithVueTs(
  pluginVue.configs['flat/essential'],
  js.configs.recommended,
  vueTsConfigs.recommended,
  configPrettier,
  {
    ignores: ['node_modules/**', 'dist/**', 'dist-ssr/**', 'coverage/**'],
    languageOptions: {
      ecmaVersion: 'latest'
    },
    rules: {
      '@typescript-eslint/no-explicit-any': 'off',
      '@typescript-eslint/no-unused-vars': 'off'
    }
  }
)
