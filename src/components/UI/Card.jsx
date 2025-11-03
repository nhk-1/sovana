import PropTypes from 'prop-types';
import { motion } from 'framer-motion';

/**
 * Elevated card with hover animations.
 */
const Card = ({ icon, title, items, index = 0 }) => (
  <motion.article
    initial={{ opacity: 0, y: 24 }}
    whileInView={{ opacity: 1, y: 0 }}
    viewport={{ once: true, amount: 0.35 }}
    transition={{ duration: 0.6, ease: 'easeOut', delay: index * 0.08 }}
    whileHover={{ translateY: -6 }}
    className="flex h-full flex-col rounded-2xl border border-brand/10 bg-white/90 p-6 shadow-subtle backdrop-blur dark:border-white/5 dark:bg-white/10"
  >
    <div className="mb-4 flex items-center gap-3 text-brand dark:text-white">
      <span className="flex h-12 w-12 items-center justify-center rounded-full bg-accent/20 text-2xl text-brand">{icon}</span>
      <h3 className="text-lg font-semibold text-brand dark:text-white">{title}</h3>
    </div>
    <ul className="mt-auto space-y-2 text-sm text-brand/80 dark:text-white/80">
      {items.map((item) => (
        <li key={item} className="flex items-start gap-2">
          <span aria-hidden="true" className="mt-1 inline-block h-1.5 w-1.5 rounded-full bg-accent"></span>
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
  index: PropTypes.number,
};

export default Card;
