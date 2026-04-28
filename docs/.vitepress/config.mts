import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: 'obacht',
  description: 'Security configuration scanner for developer environments',
	lang: "en-US",
	lastUpdated: true,
	appearance: "dark",
	ignoreDeadLinks: false,
  base: '/obacht/',
	sitemap: {
		hostname: 'https://foomo.github.io/obacht',
	},
  themeConfig: {
		// https://vitepress.dev/reference/default-theme-config
		logo: '/logo.png',
		outline: [2, 4],
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
			{
				text: 'Contributing',
				collapsed: true,
				items: [
					{
						text: "Guideline",
						link: '/CONTRIBUTING.md',
					},
					{
						text: "Code of conduct",
						link: '/CODE_OF_CONDUCT.md',
					},
					{
						text: "Security guidelines",
						link: '/SECURITY.md',
					},
				],
			},
		],
    socialLinks: [
      { icon: 'github', link: 'https://github.com/foomo/obacht' },
    ],
		editLink: {
			pattern: 'https://github.com/foomo/obacht/edit/main/docs/:path',
		},
		search: {
			provider: 'local',
		},
		footer: {
			message: 'Made with ♥ <a href="https://www.foomo.org">foomo</a> by <a href="https://www.bestbytes.com">bestbytes</a>',
		},
  },
	markdown: {
		// https://github.com/vuejs/vitepress/discussions/3724
		theme: {
			light: 'catppuccin-latte',
			dark: 'catppuccin-frappe',
		}
	},
	head: [
		['meta', { name: 'theme-color', content: '#ffffff' }],
		['link', { rel: 'icon', href: '/logo.png' }],
		['meta', { name: 'author', content: 'foomo by bestbytes' }],
		// OpenGraph
		['meta', { property: 'og:title', content: 'foomo/obacht' }],
		[
			'meta',
			{
				property: 'og:image',
				content: 'https://github.com/foomo/obacht/blob/main/docs/public/banner.png?raw=true',
			},
		],
		[
			'meta',
			{
				property: 'og:description',
				content: 'Security configuration scanner for developer environments',
			},
		],
		['meta', { name: 'twitter:card', content: 'summary_large_image' }],
		[
			'meta',
			{
				name: 'twitter:image',
				content: 'https://github.com/foomo/obacht/blob/main/docs/public/banner.png?raw=true',
			},
		],
		[
			'meta', { name: 'viewport', content: 'width=device-width, initial-scale=1.0, viewport-fit=cover',
		},
		],
	]
})
