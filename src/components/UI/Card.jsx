import PropTypes from 'prop-types';
import { motion } from 'framer-motion';

/**
 * Elevated card with hover animations.
 */
const Card = ({ icon, title, items }) => (
  <motion.article
    whileHover={{ translateY: -6 }}
    className="flex h-full flex-col rounded-2xl border border-ocean/10 bg-white/70 p-6 shadow-subtle backdrop-blur dark:bg-white/5"
  >
    <div className="mb-4 flex items-center gap-3 text-ocean dark:text-white">
      <span className="flex h-12 w-12 items-center justify-center rounded-full bg-ocean/10 text-2xl">{icon}</span>
      <h3 className="text-lg font-semibold text-midnight dark:text-white">{title}</h3>
    </div>
    <ul className="mt-auto space-y-2 text-sm text-midnight/80 dark:text-white/80">
      {items.map((item) => (
        <li key={item} className="flex items-start gap-2">
          <span aria-hidden="true" className="mt-1 inline-block h-1.5 w-1.5 rounded-full bg-ocean"></span>
          <span>{item}</span>
        </li>
      ))}
    </ul>
  </motion.article>
);

Card.propTypes = {
  icon: PropTypes.node.isRequired,
  title: PropTypes.string.isRequired,
  items: PropTypes.arrayOf(PropTypes.string).isRequired,
};

export default Card;
