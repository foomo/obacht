import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Bouncer',
  description: 'Security configuration scanner for developer environments',
  base: '/bouncer/',
  themeConfig: {
    nav: [
      { text: 'Guide', link: '/guide/getting-started' },
      { text: 'Rules', link: '/rules/' },
      { text: 'Architecture', link: '/architecture' },
    ],
    sidebar: [
      {
        text: 'Guide',
        items: [
          { text: 'Getting Started', link: '/guide/getting-started' },
          { text: 'Usage', link: '/guide/usage' },
          { text: 'Custom Rules', link: '/guide/custom-rules' },
        ],
      },
      {
        text: 'Rules',
        items: [
          { text: 'Overview', link: '/rules/' },
          { text: 'Credentials', link: '/rules/credentials' },
          { text: 'Docker', link: '/rules/docker' },
          { text: 'Environment', link: '/rules/env' },
          { text: 'Git', link: '/rules/git' },
          { text: 'Kubernetes', link: '/rules/kube' },
          { text: 'OS', link: '/rules/os' },
          { text: 'PATH', link: '/rules/path' },
          { text: 'Privacy', link: '/rules/privacy' },
          { text: 'Shell', link: '/rules/shell' },
          { text: 'SSH', link: '/rules/ssh' },
          { text: 'Tools', link: '/rules/tools' },
        ],
      },
      {
        text: 'Architecture',
        link: '/architecture',
      },
    ],
    socialLinks: [
      { icon: 'github', link: 'https://github.com/franklinkim/bouncer' },
    ],
  },
})
