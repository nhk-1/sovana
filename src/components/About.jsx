import Section from './UI/Section.jsx';

const About = () => (
  <Section id="apropos" title="À propos" description="Une équipe engagée">
    <div className="mx-auto max-w-3xl space-y-4 text-center text-lg text-midnight/80 dark:text-white/80">
      <p>
        ESN à taille humaine, Sovana place la collaboration au cœur de chaque mission. Nous co-construisons avec nos clients
        des trajectoires numériques ambitieuses, en privilégiant une approche agile, orientée résultats et durable.
      </p>
      <p>
        Nos experts accompagnent l’ensemble du cycle de vie de vos produits : cadrage, conception, développement,
        industrialisation et exploitation continue. Nous alignons nos recommandations sur vos objectifs métiers pour générer de
        la valeur mesurable.
      </p>
    </div>
  </Section>
);

export default About;
