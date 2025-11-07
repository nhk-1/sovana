import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import Section from './UI/Section.jsx';
import Button from './UI/Button.jsx';

const Contact = () => {
  const [isSent, setIsSent] = useState(false);

  // Handles the mock submission and displays a temporary toast for accessibility demo purposes.
  const handleSubmit = (event) => {
    event.preventDefault();
    setIsSent(true);
    event.target.reset();
    setTimeout(() => setIsSent(false), 4000);
  };

  return (
    <Section id="contact" title="Contact" description="Discutons de votre projet">
      <div className="mx-auto max-w-2xl">
        <form onSubmit={handleSubmit} className="space-y-6" aria-label="Formulaire de contact">
          <div>
            <label htmlFor="name" className="mb-2 block text-sm font-semibold text-brand dark:text-white">
              Nom
            </label>
            <input
              id="name"
              name="name"
              type="text"
              required
              className="w-full rounded-xl border border-brand/20 bg-white px-4 py-3 text-brand shadow-subtle placeholder:text-brand/40 focus:border-accent focus:outline-none focus:ring-2 focus:ring-accent/40 dark:border-white/20 dark:bg-white/10 dark:text-white"
              placeholder="Votre nom"
              autoComplete="name"
            />
          </div>
          <div>
            <label htmlFor="email" className="mb-2 block text-sm font-semibold text-brand dark:text-white">
              Email
            </label>
            <input
              id="email"
              name="email"
              type="email"
              required
              className="w-full rounded-xl border border-brand/20 bg-white px-4 py-3 text-brand shadow-subtle placeholder:text-brand/40 focus:border-accent focus:outline-none focus:ring-2 focus:ring-accent/40 dark:border-white/20 dark:bg-white/10 dark:text-white"
              placeholder="vous@entreprise.com"
              autoComplete="email"
            />
          </div>
          <div>
            <label htmlFor="message" className="mb-2 block text-sm font-semibold text-brand dark:text-white">
              Message
            </label>
            <textarea
              id="message"
              name="message"
              required
              rows="5"
              className="w-full rounded-xl border border-brand/20 bg-white px-4 py-3 text-brand shadow-subtle placeholder:text-brand/40 focus:border-accent focus:outline-none focus:ring-2 focus:ring-accent/40 dark:border-white/20 dark:bg-white/10 dark:text-white"
              placeholder="Parlez-nous de vos enjeux..."
            ></textarea>
          </div>
          <div className="flex flex-col items-start gap-2 sm:flex-row sm:items-center">
            <Button>Envoyer</Button>
            <span className="text-sm text-brand/70 dark:text-white/60">Email direct : contact@sovana.fr</span>
          </div>
        </form>
        <AnimatePresence>
          {isSent && (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              role="status"
              aria-live="polite"
              className="mt-6 rounded-xl bg-brand px-4 py-3 text-sm font-medium text-white shadow-subtle"
            >
              Message envoyé (démo)
            </motion.div>
          )}
        </AnimatePresence>
      </div>
    </Section>
  );
};

export default Contact;
