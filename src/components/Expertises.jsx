import { motion } from 'framer-motion';
import Section from './UI/Section.jsx';

const TECHS = [
  { name: 'Docker', icon: 'docker' },
  { name: 'Linux', icon: 'linux' },
  { name: 'Microsoft Azure', icon: 'azure' },
  { name: 'Amazon Web Services', icon: 'aws' },
  { name: 'Terraform', icon: 'terraform' },
  { name: 'Ansible', icon: 'ansible' },
  { name: 'React', icon: 'react' },
  { name: 'Node.js', icon: 'node' },
  { name: 'Python', icon: 'python' },
  { name: 'PostgreSQL', icon: 'postgresql' },
];

const getIcon = (icon) => new URL(`../assets/tech/${icon}.svg`, import.meta.url).href;

const Expertises = () => (
  <Section id="expertises" title="Nos expertises" description="Stack éprouvée">
    <div className="grid gap-6 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-5">
      {TECHS.map((tech, index) => (
        <motion.div
          key={tech.name}
          initial={{ opacity: 0, scale: 0.9 }}
          whileInView={{ opacity: 1, scale: 1 }}
          viewport={{ once: true, amount: 0.3 }}
          transition={{ duration: 0.5, ease: 'easeOut', delay: index * 0.04 }}
          whileHover={{ translateY: -4 }}
          className="group relative flex flex-col items-center gap-3 rounded-2xl border border-brand/10 bg-white/90 p-6 text-center shadow-subtle transition hover:border-accent/40 focus-within:border-accent/40 dark:border-white/10 dark:bg-white/10"
        >
          <img src={getIcon(tech.icon)} alt="" role="presentation" className="h-16 w-16" loading="lazy" />
          <button type="button" className="text-sm font-semibold text-brand dark:text-white" aria-describedby={`${tech.icon}-tooltip`}>
            {tech.name}
          </button>
          <span
            id={`${tech.icon}-tooltip`}
            role="tooltip"
            className="pointer-events-none absolute bottom-2 left-1/2 w-max -translate-x-1/2 translate-y-full rounded-lg bg-brand px-3 py-1 text-xs text-white opacity-0 transition group-hover:translate-y-2 group-hover:opacity-100 group-focus-within:translate-y-2 group-focus-within:opacity-100"
          >
            Expertise certifiée
          </span>
        </motion.div>
      ))}
    </div>
  </Section>
);

export default Expertises;
