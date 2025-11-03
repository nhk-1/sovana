import Container from './UI/Container.jsx';

const Footer = () => (
  <footer className="border-t border-white/10 bg-white/80 py-8 dark:bg-midnight/80" role="contentinfo">
    <Container className="flex flex-col gap-4 text-sm text-midnight/70 dark:text-white/70 md:flex-row md:items-center md:justify-between">
      <div className="flex flex-col gap-2">
        <span className="font-semibold text-midnight dark:text-white">© 2025 Sovana</span>
        <a href="#" className="hover:text-ocean">
          Mentions légales
        </a>
      </div>
      <div className="flex items-center gap-6">
        <a href="https://www.linkedin.com" className="hover:text-ocean" aria-label="LinkedIn Sovana">
          LinkedIn
        </a>
        <a href="https://github.com" className="hover:text-ocean" aria-label="GitHub Sovana">
          GitHub
        </a>
      </div>
    </Container>
  </footer>
);

export default Footer;
