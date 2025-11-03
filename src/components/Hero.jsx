import { motion } from 'framer-motion';
import Button from './UI/Button.jsx';
import Container from './UI/Container.jsx';
import Badge from './UI/Badge.jsx';

const Hero = () => (
  <section
    id="accueil"
    className="relative overflow-hidden bg-gradient-to-b from-white via-white to-white py-24 dark:from-night dark:via-night dark:to-night"
    aria-labelledby="hero-title"
  >
    <div className="absolute inset-0 -z-10 bg-[radial-gradient(circle_at_top,_rgba(14,186,227,0.3),_transparent_60%)]" aria-hidden="true" />
    <Container className="flex flex-col items-start gap-10 text-left lg:flex-row lg:items-center lg:justify-between">
      <motion.div
        initial={{ opacity: 0, y: 30 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.8, ease: 'easeOut' }}
        className="max-w-xl space-y-6"
      >
        <Badge>ESN française</Badge>
        <h1 id="hero-title" className="text-4xl font-bold leading-tight text-brand dark:text-white sm:text-5xl">
          Sovana – L’ESN qui propulse vos projets numériques.
        </h1>
        <p className="text-lg text-brand/80 dark:text-white/80">
          Conseil, développement et intégration sur mesure pour les entreprises innovantes. Notre équipe pluridisciplinaire
          délivre des solutions fiables et sécurisées, alignées sur vos objectifs business.
        </p>
        <div className="flex flex-col gap-4 sm:flex-row">
          <Button href="#services">Découvrir nos services</Button>
          <Button href="mailto:contact@sovana.fr" variant="secondary">
            contact@sovana.fr
          </Button>
        </div>
      </motion.div>
      <motion.div
        initial={{ opacity: 0, x: 50 }}
        animate={{ opacity: 1, x: 0 }}
        transition={{ duration: 0.9, delay: 0.2, ease: 'easeOut' }}
        className="relative w-full max-w-md"
        role="presentation"
        aria-hidden="true"
      >
        <div className="aspect-square rounded-full bg-accent/25" />
        <div className="absolute inset-6 rounded-[2.5rem] border border-dashed border-accent/60" />
        <div className="absolute inset-12 rounded-3xl bg-gradient-to-br from-brand to-accent" />
      </motion.div>
    </Container>
  </section>
);

export default Hero;
