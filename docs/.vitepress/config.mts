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
          { text: 'SSH', link: '/rules/ssh' },
          { text: 'Git', link: '/rules/git' },
          { text: 'Docker', link: '/rules/docker' },
          { text: 'Kubernetes', link: '/rules/kube' },
          { text: 'Environment', link: '/rules/env' },
          { text: 'Shell', link: '/rules/shell' },
          { text: 'Tools', link: '/rules/tools' },
          { text: 'PATH', link: '/rules/path' },
          { text: 'OS', link: '/rules/os' },
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
