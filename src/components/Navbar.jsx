import { useState } from 'react';
import { Link } from 'react-scroll';
import Button from './UI/Button.jsx';
import Container from './UI/Container.jsx';
import logo from '../assets/logo.svg';

const NAV_LINKS = [
  { id: 'accueil', label: 'Accueil' },
  { id: 'services', label: 'Services' },
  { id: 'expertises', label: 'Expertises' },
  { id: 'apropos', label: 'À propos' },
  { id: 'contact', label: 'Contact' },
];

const Navbar = () => {
  const [isOpen, setIsOpen] = useState(false);

  const toggleMenu = () => setIsOpen((prev) => !prev);

  return (
    <header className="sticky top-0 z-50 border-b border-brand/10 bg-canvas/85 backdrop-blur dark:border-night/40 dark:bg-night/85" role="banner">
      <Container className="flex items-center justify-between py-4">
        <a href="#accueil" className="flex items-center gap-3 text-brand" aria-label="Retour à l’accueil">
          <img src={logo} alt="Sovana" className="h-10 w-auto" loading="lazy" />
          <span className="text-xl font-semibold text-brand dark:text-white">Sovana</span>
        </a>
        <nav aria-label="Navigation principale" className="hidden items-center gap-8 md:flex">
          {NAV_LINKS.map((link) => (
            <Link
              key={link.id}
              to={link.id}
              smooth
              duration={600}
              offset={-80}
              className="cursor-pointer text-sm font-medium text-brand/75 transition hover:text-accent focus:outline-none focus-visible:text-accent dark:text-white/80"
            >
              {link.label}
            </Link>
          ))}
          <Button href="#contact" variant="primary">
            Nous contacter
          </Button>
        </nav>
        <button
          type="button"
          className="text-brand md:hidden"
          aria-expanded={isOpen}
          aria-controls="mobile-menu"
          onClick={toggleMenu}
        >
          <span className="sr-only">{isOpen ? 'Fermer le menu' : 'Ouvrir le menu'}</span>
          <div className="space-y-1.5">
            {[0, 1, 2].map((line) => (
              <span
                key={line}
                className="block h-0.5 w-6 bg-brand transition dark:bg-white"
              ></span>
            ))}
          </div>
        </button>
      </Container>
      {isOpen && (
        <nav
          id="mobile-menu"
          aria-label="Navigation principale mobile"
          className="md:hidden"
        >
          <Container className="flex flex-col gap-4 pb-6">
            {NAV_LINKS.map((link) => (
              <Link
                key={link.id}
                to={link.id}
                smooth
                duration={600}
                offset={-80}
                onClick={() => setIsOpen(false)}
                className="cursor-pointer text-base font-medium text-brand/75 transition hover:text-accent focus:outline-none focus-visible:text-accent dark:text-white/80"
              >
                {link.label}
              </Link>
            ))}
            <Button href="#contact" onClick={() => setIsOpen(false)}>
              Nous contacter
            </Button>
          </Container>
        </nav>
      )}
    </header>
  );
};

export default Navbar;
