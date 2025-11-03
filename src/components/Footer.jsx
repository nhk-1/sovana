import Container from './UI/Container.jsx';

const Footer = () => (
  <footer className="border-t border-brand/10 bg-canvas py-8 dark:border-night/40 dark:bg-night" role="contentinfo">
    <Container className="flex flex-col gap-4 text-sm text-brand/70 dark:text-white/70 md:flex-row md:items-center md:justify-between">
      <div className="flex flex-col gap-2">
        <span className="font-semibold text-brand dark:text-white">© 2025 Sovana</span>
        <a href="#" className="transition hover:text-accent">
          Mentions légales
        </a>
      </div>
      <div className="flex items-center gap-6">
        <a href="https://www.linkedin.com" className="transition hover:text-accent" aria-label="LinkedIn Sovana">
          LinkedIn
        </a>
        <a href="https://github.com" className="transition hover:text-accent" aria-label="GitHub Sovana">
          GitHub
        </a>
      </div>
    </Container>
  </footer>
);

export default Footer;
