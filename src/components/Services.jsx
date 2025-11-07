import Section from './UI/Section.jsx';
import Card from './UI/Card.jsx';

const SERVICES = [
  {
    icon: '🧠',
    title: 'Développement sur mesure',
    items: ['Applications web & mobile', 'Architecture modulaire', 'Intégrations API', 'Pilotage produit agile'],
  },
  {
    icon: '☁️',
    title: 'Infrastructure & Cloud',
    items: ['Modernisation cloud', 'Infrastructure as Code', 'Observabilité & monitoring', 'Optimisation des coûts'],
  },
  {
    icon: '🛡️',
    title: 'Cybersécurité & Réseaux',
    items: ['Audits & pentests', 'Zero Trust & IAM', 'Sécurisation des flux', 'Surveillance proactive'],
  },
  {
    icon: '🤝',
    title: 'Support & Infogérance',
    items: ['SLA personnalisés', 'Support 24/7', 'MCO / MCS', 'Gestion des incidents'],
  },
];

const Services = () => (
  <Section
    id="services"
    title="Nos services"
    description="Solutions end-to-end"
  >
    <div className="grid gap-6 md:grid-cols-2 xl:grid-cols-4">
      {SERVICES.map((service, index) => (
        <Card key={service.title} icon={service.icon} title={service.title} items={service.items} index={index} />
      ))}
    </div>
  </Section>
);

export default Services;
