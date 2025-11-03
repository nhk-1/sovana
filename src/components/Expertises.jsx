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
      {TECHS.map((tech) => (
        <div
          key={tech.name}
          className="group relative flex flex-col items-center gap-3 rounded-2xl border border-ocean/10 bg-white/70 p-6 text-center shadow-subtle transition hover:border-ocean/30 focus-within:border-ocean/30 dark:bg-white/5"
        >
          <img src={getIcon(tech.icon)} alt="" role="presentation" className="h-16 w-16" loading="lazy" />
          <button type="button" className="text-sm font-semibold text-midnight dark:text-white" aria-describedby={`${tech.icon}-tooltip`}>
            {tech.name}
          </button>
          <span
            id={`${tech.icon}-tooltip`}
            role="tooltip"
            className="pointer-events-none absolute bottom-2 left-1/2 w-max -translate-x-1/2 translate-y-full rounded-lg bg-midnight px-3 py-1 text-xs text-white opacity-0 transition group-hover:translate-y-2 group-hover:opacity-100 group-focus-within:translate-y-2 group-focus-within:opacity-100"
          >
            Expertise certifiée
          </span>
        </div>
      ))}
    </div>
  </Section>
);

export default Expertises;
