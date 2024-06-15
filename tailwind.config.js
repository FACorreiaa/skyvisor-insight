/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./app/**/*.{css,view,js,templ,html}'],
  theme: {
    container: {
      center: true,
      padding: "2rem",
      screens: {
        "2xl": "1400px"
      }
    },
    extend: {
      colors: {
        border: "hsl(var(--border) / <alpha-value>)",
        input: "hsl(var(--input) / <alpha-value>)",
        ring: "hsl(var(--ring) / <alpha-value>)",
        background: "hsl(var(--background) / <alpha-value>)",
        foreground: "hsl(var(--foreground) / <alpha-value>)",
        primary: {
          DEFAULT: "hsl(var(--primary) / <alpha-value>)",
          foreground: "hsl(var(--primary-foreground) / <alpha-value>)"
        },
        secondary: {
          DEFAULT: "hsl(var(--secondary) / <alpha-value>)",
          foreground: "hsl(var(--secondary-foreground) / <alpha-value>)"
        },
        destructive: {
          DEFAULT: "hsl(var(--destructive) / <alpha-value>)",
          foreground: "hsl(var(--destructive-foreground) / <alpha-value>)"
        },
        muted: {
          DEFAULT: "hsl(var(--muted) / <alpha-value>)",
          foreground: "hsl(var(--muted-foreground) / <alpha-value>)"
        },
        accent: {
          DEFAULT: "hsl(var(--accent) / <alpha-value>)",
          foreground: "hsl(var(--accent-foreground) / <alpha-value>)"
        },
        popover: {
          DEFAULT: "hsl(var(--popover) / <alpha-value>)",
          foreground: "hsl(var(--popover-foreground) / <alpha-value>)"
        },
        card: {
          DEFAULT: "hsl(var(--card) / <alpha-value>)",
          foreground: "hsl(var(--card-foreground) / <alpha-value>)"
        }
      },
      borderRadius: {
        lg: "var(--radius)",
        md: "calc(var(--radius) - 2px)",
        sm: "calc(var(--radius) - 4px)"
      },
      spacing: {
        '8xl': '96rem',
        '9xl': '128rem',
      },
      borderRadius: {
        '4xl': '2rem',
      },
    },
    colors: {
      blue: '#1fb6ff',
      purple: '#7e5bef',
      pink: '#ff49db',
      orange: '#ff7849',
      green: '#13ce66',
      yellow: '#ffc82c',
      'gray-dark': '#273444',
      gray: '#8492a6',
      'gray-light': '#d3dce6',
    },
    fontFamily: {
      sans: ['Graphik', 'sans-serif'],
      serif: ['Merriweather', 'serif'],
      lato: ['Lato', 'sans-serif'],
    },
  },
  daisyui: {
    themes: [
      'light',
      'dark',
      'nord',
      'cyberpunk',
      'pastel',
      'cupcake',
      'night',
      {
        'catppuccin-latte': {
          primary: '#1e66f5',
          secondary: '#ea76cb',
          accent: '#179299',
          neutral: '#dce0e8',
          'base-100': '#eff1f5',
          info: '#209fb5',
          success: '#40a02b',
          warning: '#df8e1d',
          error: '#d20f39',
        },
        'catppuccin-frappe': {
          primary: '#8caaee',
          secondary: '#f4b8e4',
          accent: '#81c8be',
          neutral: '#232634',
          'base-100': '#303446',
          info: '#85c1dc',
          success: '#a6d189',
          warning: '#e5c890',
          error: '#e78284',
        },
        'catppuccin-macchiato': {
          primary: '#8aadf4',
          secondary: '#f5bde6',
          accent: '#8bd5ca',
          neutral: '#181926',
          'base-100': '#24273a',
          info: '#7dc4e4',
          success: '#a6da95',
          warning: '#eed49f',
          error: '#ed8796',
        },
        'catppuccin-mocha': {
          primary: '#89b4fa',
          secondary: '#f5c2e7',
          accent: '#94e2d5',
          neutral: '#11111b',
          'base-100': '#1e1e2e',
          info: '#74c7ec',
          success: '#a6e3a1',
          warning: '#f9e2af',
          error: '#f38ba8',
        },
      },
    ],
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
    require('daisyui'),
    require('autoprefixer'),
  ],
};
