import { useRef } from 'react';
import PropTypes from 'prop-types';
import { motion, useInView } from 'framer-motion';
import Container from './Container.jsx';

/**
 * Semantic section wrapper providing spacing and optional title.
 */
const Section = ({ id, title, description, children, className = '', as = 'section' }) => {
  const Component = as;
  const labelledBy = title ? `${id}-title` : undefined;
  const ref = useRef(null);
  const isInView = useInView(ref, { once: true, margin: '-120px' });

  return (
    <Component id={id} aria-labelledby={labelledBy} className={className}>
      <motion.div
        ref={ref}
        initial={{ opacity: 0, y: 48 }}
        animate={isInView ? { opacity: 1, y: 0 } : {}}
        transition={{ duration: 0.8, ease: 'easeOut' }}
        className="w-full"
      >
        <Container className="py-16 sm:py-20">
          {title ? (
            <div className="mb-10 flex flex-col gap-3 text-center">
              {description && (
                <p className="text-sm uppercase tracking-[0.2em] text-accent/80 dark:text-accent/70">{description}</p>
              )}
              <h2 id={labelledBy} className="text-3xl font-bold text-brand dark:text-white sm:text-4xl">
                {title}
              </h2>
            </div>
          ) : null}
          {children}
        </Container>
      </motion.div>
    </Component>
  );
};

Section.propTypes = {
  id: PropTypes.string,
  title: PropTypes.string,
  description: PropTypes.string,
  children: PropTypes.node.isRequired,
  className: PropTypes.string,
  as: PropTypes.oneOf(['section', 'div']),
};

export default Section;
