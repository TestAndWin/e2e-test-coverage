import pluginVue from 'eslint-plugin-vue'
import { defineConfigWithVueTs, vueTsConfigs } from '@vue/eslint-config-typescript'
import configPrettier from '@vue/eslint-config-prettier/skip-formatting'
import js from '@eslint/js'

export default defineConfigWithVueTs(
  // Must stay a config object with only "ignores" - that is what makes it a
  // global ignore. Combined with other keys it would only apply to that block.
  { ignores: ['node_modules/**', 'dist/**', 'dist-ssr/**', 'coverage/**'] },
  pluginVue.configs['flat/essential'],
  js.configs.recommended,
  vueTsConfigs.recommended,
  configPrettier,
  {
    languageOptions: {
      ecmaVersion: 'latest'
    },
    rules: {
      '@typescript-eslint/no-explicit-any': 'off',
      '@typescript-eslint/no-unused-vars': 'off'
    }
  }
)
